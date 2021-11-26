package message

import (
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"log"
	mathRand "math/rand"
)


type PacketType uint8
const(
	SYN 		PacketType = 1 + iota //공개키 주고 받음
	SYN_ACK			//패킷 코드, 키 주고 받음 (XOR 용), 이 패킷은 RSA 암호화가 기본
	ESTABLISHED 	//이 패킷은 연결된 후 패킷, 여기서부터는 서로 공개키와 패킷 코드, 키를 알고 있으므로 인코딩은 선택하면됨
)

type CryptoType uint8
const (
	XOR CryptoType = 1 + iota
	RSA
)

type Header struct{
	packetType		PacketType
	cryptoType	 	CryptoType
	randKey			uint8
	payloadLength	uint16
	checkSum		uint8
}

const (
	headerSize  = 6
	payloadSize = 300
	bufferSize  = headerSize + payloadSize
)

type Message struct{
	buffer []byte

	front uint32
	rear  uint32

	order binary.ByteOrder
}

func NewMessage(isLittleEndian bool) *Message {
	var msg = Message{
		buffer: make([]byte, bufferSize),

		front: headerSize,
		rear:  headerSize,

		order: binary.LittleEndian,
	}
	if isLittleEndian == false {
		msg.order = binary.BigEndian
	}

	return &msg
}

func (msg *Message) GetBuffer() []byte{
	return msg.buffer
}

func (msg *Message) GetHeaderBuffer() []byte{
	return msg.buffer[:headerSize]
}

func (msg *Message) GetPayloadBuffer() []byte{
	return msg.buffer[headerSize:]
}

func (msg *Message) GetPayloadLength() uint32{
	return msg.rear - msg.front
}

func (msg *Message) GetLength() int {
	return len(msg.buffer)
}

func (msg *Message) SetHeader(packetType PacketType, cryptoType CryptoType){
	var header Header
	header.packetType = packetType
	header.cryptoType = cryptoType
	header.randKey = uint8(mathRand.Intn(256))
	header.payloadLength = uint16(msg.GetPayloadLength())
	header.checkSum = msg.generateChecksum()

	msg.packHeader(header) //pragma pack(1)
}


func (msg *Message) Push(value interface{}) uint32 {
	var pushSize uint32
	var tmpBuffer []byte
	tmpBuffer = make([]byte, bufferSize)

	switch value.(type){
	case bool, byte:
		tmpBuffer[0] = value.(byte)
		pushSize = 1
		break
	case uint16, int16:
		msg.order.PutUint16(tmpBuffer, value.(uint16))
		pushSize = 2
		break
	case uint32, int32:
		msg.order.PutUint32(tmpBuffer, value.(uint32))
		pushSize = 4
		break
	case uint64, int64:
		msg.order.PutUint64(tmpBuffer, value.(uint64))
		pushSize = 8
		break
	case string:
		length := uint16(len(value.(string)))
		msg.Push(length)
		pushSize = uint32(copy(tmpBuffer, value.(string)))
		break
	case []byte:
		length := uint16(len(value.([]byte)))
		msg.Push(length)
		pushSize = uint32(copy(tmpBuffer, value.([]byte)))
		break
	default:
		return 0
	}

	if msg.getFreeLength() < pushSize {
		return 0
	}

	copy(msg.buffer[msg.rear:], tmpBuffer)
	msg.rear += pushSize

	return pushSize
}

func (msg *Message) Peek(out_value interface{}) uint32{
	var peekSize uint32
	var tmpBuffer []byte

	switch out_value.(type){
	case *bool, *byte:
		pOutValue := out_value.(*byte)
		*pOutValue = msg.buffer[msg.front]
		peekSize = 1
		break
	case *uint16, *int16:
		tmpBuffer = msg.buffer[msg.front : msg.front+ 2]
		pOutValue := out_value.(*uint16)
		*pOutValue = msg.order.Uint16(tmpBuffer)
		peekSize = 2
		break
	case *uint32, *int32:
		tmpBuffer = msg.buffer[msg.front : msg.front+ 4]
		pOutValue := out_value.(*uint32)
		*pOutValue = msg.order.Uint32(tmpBuffer)
		peekSize = 4
		break
	case *uint64, *int64:
		tmpBuffer = msg.buffer[msg.front : msg.front+ 8]
		pOutValue := out_value.(*uint64)
		*pOutValue = msg.order.Uint64(tmpBuffer)
		peekSize = 8
		break
	case *string:
		var length uint16
		msg.Pop(&length)
		tmpBuffer = msg.buffer[msg.front : msg.front+ uint32(length)]
		pOutValue := out_value.(*string)
		*pOutValue = string(tmpBuffer)
		peekSize = uint32(length)
		break
	case *[]byte:
		var length uint16
		msg.Pop(&length)
		tmpBuffer = msg.buffer[msg.front : msg.front+ uint32(length)]
		pOutValue := out_value.(*[]byte)
		*pOutValue = tmpBuffer
		peekSize = uint32(length)
		break
	default:
		return 0
	}

	return peekSize
}

func (msg *Message) Pop(out_value interface{}) uint32 {
	popSize := msg.Peek(out_value)

	msg.front += popSize

	return popSize
}

func (msg *Message) EncodeXOR(key uint8){
	if msg.isCryptoType(XOR) != true {
		//TODO_MSG :: 로그 삽입
		return
	}

	randKey := msg.buffer[2]
	dstBuffer := msg.buffer[headerSize- 1 : msg.rear]

	num := uint32(1)
	for i := range dstBuffer {
		p := dstBuffer[i] ^ uint8(uint32(randKey) + num)
		dstBuffer[i] = p ^ uint8(uint32(key) + num)
	}
}

func (msg *Message) DecodeXOR(key uint8){
	if msg.isCryptoType(XOR) != true {
		//TODO_MSG :: 로그 삽입
		return
	}

	randKey := msg.buffer[2]
	dstBuffer := msg.buffer[headerSize- 1 : msg.rear]

	num := uint32(1)
	for i := range dstBuffer {
		p := dstBuffer[i] ^ uint8(uint32(key) + num)
		dstBuffer[i] = p ^ uint8(uint32(randKey) + num)
	}

	//체크섬 확인
	recvChecksum := dstBuffer[0]
	generatedChecksum := msg.generateChecksum()
	if recvChecksum != generatedChecksum {
		//TODO_MSG :: 로그 삽입
		log.Fatalln("mismatch check sum : ", recvChecksum, generatedChecksum)
		return;
	}
}

func (msg *Message) EncodeRSA(clntPublicKey *rsa.PublicKey){
	if msg.isCryptoType(RSA) != true {
		//TODO_MSG :: 로그 삽입
		return
	}

	cipherMsg, err := rsa.EncryptPKCS1v15(cryptoRand.Reader, clntPublicKey, msg.buffer)
	if err != nil{
		//TODO_MSG :: 로그 추가 필요
		return
	}

	cipherMsgLength := len(cipherMsg)

	msg.clear()
	copy(msg.buffer, cipherMsg)
	msg.rear += uint32(cipherMsgLength)
}

func (msg *Message) DecodeRSA(servPrivateKey *rsa.PrivateKey){
	if msg.isCryptoType(RSA) != true {
		//TODO_MSG :: 로그 삽입
		return
	}

	plainMsg, err := rsa.DecryptPKCS1v15(cryptoRand.Reader, servPrivateKey, msg.buffer)
	if err != nil {
		//TODO_MSG :: 로그 추가 필요
		return
	}

	plainMsgLength := len(plainMsg)

	msg.clear()
	copy(msg.buffer, plainMsg)
	msg.rear += uint32(plainMsgLength)
}

func (msg *Message) clear() {
	for i := range msg.buffer {
		msg.buffer[i] = 0
	}

	msg.front = headerSize
	msg.rear = headerSize
}

func (msg *Message) getFreeLength() uint32 {
	tempFront := msg.front
	tempRear := msg.rear

	return (payloadSize - 1) - (tempRear - tempFront)
}

func (msg *Message) generateChecksum() uint8{
	var total uint32

	payload := msg.buffer[msg.front:msg.rear]

	for i := range payload{
		total += uint32(payload[i])
	}

	return uint8(total % 256)
}

func (msg *Message) packHeader(header Header) {
	msg.buffer[0] = byte(header.packetType)
	msg.buffer[1] = byte(header.cryptoType)
	msg.buffer[2] = header.randKey
	msg.order.PutUint16(msg.buffer[3:5], header.payloadLength)
	msg.buffer[5] = header.checkSum
}

func (msg *Message) isCryptoType(cryptoType CryptoType) bool {
	return CryptoType(msg.buffer[1]) == cryptoType
}