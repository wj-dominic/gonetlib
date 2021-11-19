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
	HEADER_SIZE		= 6
	PAYLOAD_SIZE	= 300
	MAX_SIZE		= HEADER_SIZE + PAYLOAD_SIZE
)

type Message struct{
	Buffer	[]byte

	Front	uint32
	Rear	uint32

	Order	binary.ByteOrder
}

func NewMessage(isLittleEndian bool) *Message {
	var msg = Message{
		Buffer: make([]byte, MAX_SIZE),

		Front: HEADER_SIZE,
		Rear:  HEADER_SIZE,

		Order: binary.LittleEndian,
	}
	if isLittleEndian == false {
		msg.Order = binary.BigEndian
	}

	return &msg
}

func (msg *Message) GetBuffer() []byte{
	return msg.Buffer
}

func (msg *Message) GetHeaderBuffer() []byte{
	return msg.Buffer[:HEADER_SIZE]
}

func (msg *Message) GetPayloadBuffer() []byte{
	return msg.Buffer[HEADER_SIZE:]
}

func (msg *Message) GetPayloadLength() uint32{
	return msg.Rear - msg.Front
}

func (msg *Message) GetLength() int {
	return len(msg.Buffer)
}

func (msg *Message) SetHeader(packetType PacketType, cryptoType CryptoType){
	var header Header
	header.packetType = packetType
	header.cryptoType = cryptoType
	header.randKey = uint8(mathRand.Intn(256))
	header.payloadLength = uint16(msg.GetPayloadLength())
	header.checkSum = msg._GenerateChecksum()

	msg._PackHeader(header) //pragma pack(1)
}


func (msg *Message) Push(value interface{}) uint32 {
	var pushSize uint32
	var tmpBuffer []byte
	tmpBuffer = make([]byte, MAX_SIZE)

	switch value.(type){
	case bool, byte:
		tmpBuffer[0] = value.(byte)
		pushSize = 1
		break
	case uint16, int16:
		msg.Order.PutUint16(tmpBuffer, value.(uint16))
		pushSize = 2
		break
	case uint32, int32:
		msg.Order.PutUint32(tmpBuffer, value.(uint32))
		pushSize = 4
		break
	case uint64, int64:
		msg.Order.PutUint64(tmpBuffer, value.(uint64))
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

	if msg._GetFreeLength() < pushSize {
		return 0
	}

	copy(msg.Buffer[msg.Rear:], tmpBuffer)
	msg.Rear += pushSize

	return pushSize
}

func (msg *Message) Peek(out_value interface{}) uint32{
	var peekSize uint32
	var tmpBuffer []byte

	switch out_value.(type){
	case *bool, *byte:
		pOutValue := out_value.(*byte)
		*pOutValue = msg.Buffer[msg.Front]
		peekSize = 1
		break
	case *uint16, *int16:
		tmpBuffer = msg.Buffer[msg.Front : msg.Front + 2]
		pOutValue := out_value.(*uint16)
		*pOutValue = msg.Order.Uint16(tmpBuffer)
		peekSize = 2
		break
	case *uint32, *int32:
		tmpBuffer = msg.Buffer[msg.Front : msg.Front + 4]
		pOutValue := out_value.(*uint32)
		*pOutValue = msg.Order.Uint32(tmpBuffer)
		peekSize = 4
		break
	case *uint64, *int64:
		tmpBuffer = msg.Buffer[msg.Front : msg.Front + 8]
		pOutValue := out_value.(*uint64)
		*pOutValue = msg.Order.Uint64(tmpBuffer)
		peekSize = 8
		break
	case *string:
		var length uint16
		msg.Pop(&length)
		tmpBuffer = msg.Buffer[msg.Front : msg.Front + uint32(length)]
		pOutValue := out_value.(*string)
		*pOutValue = string(tmpBuffer)
		peekSize = uint32(length)
		break
	case *[]byte:
		var length uint16
		msg.Pop(&length)
		tmpBuffer = msg.Buffer[msg.Front : msg.Front + uint32(length)]
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

	msg.Front += popSize

	return popSize
}

func (msg *Message) EncodeXOR(key uint8){
	if msg._IsCryptoType(XOR) != true {
		//TODO_MSG :: 로그 삽입
		return
	}

	randKey := msg.Buffer[2]
	dstBuffer := msg.Buffer[HEADER_SIZE - 1 : msg.Rear]

	num := uint32(1)
	for i := range dstBuffer {
		p := dstBuffer[i] ^ uint8(uint32(randKey) + num)
		dstBuffer[i] = p ^ uint8(uint32(key) + num)
	}
}

func (msg *Message) DecodeXOR(key uint8){
	if msg._IsCryptoType(XOR) != true {
		//TODO_MSG :: 로그 삽입
		return
	}

	randKey := msg.Buffer[2]
	dstBuffer := msg.Buffer[HEADER_SIZE - 1 : msg.Rear]

	num := uint32(1)
	for i := range dstBuffer {
		p := dstBuffer[i] ^ uint8(uint32(key) + num)
		dstBuffer[i] = p ^ uint8(uint32(randKey) + num)
	}

	//체크섬 확인
	recvChecksum := dstBuffer[0]
	generatedChecksum := msg._GenerateChecksum()
	if recvChecksum != generatedChecksum {
		//TODO_MSG :: 로그 삽입
		log.Fatalln("mismatch check sum : ", recvChecksum, generatedChecksum)
		return;
	}
}

func (msg *Message) EncodeRSA(clntPublicKey *rsa.PublicKey){
	if msg._IsCryptoType(RSA) != true {
		//TODO_MSG :: 로그 삽입
		return
	}

	cipherMsg, err := rsa.EncryptPKCS1v15(cryptoRand.Reader, clntPublicKey, msg.Buffer)
	if err != nil{
		//TODO_MSG :: 로그 추가 필요
		return
	}

	cipherMsgLength := len(cipherMsg)

	msg._Clear()
	copy(msg.Buffer, cipherMsg)
	msg.Rear += uint32(cipherMsgLength)
}

func (msg *Message) DecodeRSA(servPrivateKey *rsa.PrivateKey){
	if msg._IsCryptoType(RSA) != true {
		//TODO_MSG :: 로그 삽입
		return
	}

	plainMsg, err := rsa.DecryptPKCS1v15(cryptoRand.Reader, servPrivateKey, msg.Buffer)
	if err != nil {
		//TODO_MSG :: 로그 추가 필요
		return
	}

	plainMsgLength := len(plainMsg)

	msg._Clear()
	copy(msg.Buffer, plainMsg)
	msg.Rear += uint32(plainMsgLength)
}

func (msg *Message) _Clear() {
	for i := range msg.Buffer{
		msg.Buffer[i] = 0
	}

	msg.Front = HEADER_SIZE
	msg.Rear = HEADER_SIZE
}

func (msg *Message) _GetFreeLength() uint32 {
	tempFront := msg.Front
	tempRear := msg.Rear

	return (PAYLOAD_SIZE - 1) - (tempRear - tempFront)
}

func (msg *Message) _GenerateChecksum() uint8{
	var total uint32

	payload := msg.Buffer[msg.Front:msg.Rear]

	for i := range payload{
		total += uint32(payload[i])
	}

	return uint8(total % 256)
}

func (msg *Message) _PackHeader(header Header) {
	msg.Buffer[0] = byte(header.packetType)
	msg.Buffer[1] = byte(header.cryptoType)
	msg.Buffer[2] = header.randKey
	msg.Order.PutUint16(msg.Buffer[3:5], header.payloadLength)
	msg.Buffer[5] = header.checkSum
}

func (msg *Message) _IsCryptoType(cryptoType CryptoType) bool {
	return CryptoType(msg.Buffer[1]) == cryptoType
}