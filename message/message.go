package message

import (
	"bytes"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"fmt"
	"log"
	mathRand "math/rand"
	"reflect"
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

type NetHeader struct{
	PacketType 		PacketType
	CryptoType    	CryptoType
	RandKey       	uint8
	PayloadLength 	uint16
	CheckSum      	uint8
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
	return msg.buffer[msg.front:msg.rear]
}

func (msg *Message) GetPayloadLength() uint32{
	return msg.rear - msg.front
}

func (msg *Message) GetLength() int {
	return len(msg.buffer)
}

func (msg *Message) SetHeader(packetType PacketType, cryptoType CryptoType){
	var header NetHeader
	header.PacketType = packetType
	header.CryptoType = cryptoType
	header.RandKey = uint8(mathRand.Intn(256))
	header.PayloadLength = uint16(msg.GetPayloadLength())
	header.CheckSum = msg.generateChecksum()

	msg.packHeader(header) //pragma pack(1)
}


func (msg *Message) Push(value interface{}) uint32 {
	var pushSize uint32
	var tmpBuffer []byte
	tmpBuffer = make([]byte, bufferSize)

	switch reflect.TypeOf(value).Kind() {
	case reflect.Bool, reflect.Uint8, reflect.Int8:
		tmpBuffer[0] = value.(byte)
		pushSize = 1
		break
	case reflect.Uint16, reflect.Int16:
		msg.order.PutUint16(tmpBuffer, value.(uint16))
		pushSize = 2
		break
	case reflect.Uint32, reflect.Int32:
		msg.order.PutUint32(tmpBuffer, value.(uint32))
		pushSize = 4
		break
	case reflect.Uint64, reflect.Int64:
		msg.order.PutUint64(tmpBuffer, value.(uint64))
		pushSize = 8
		break
	case reflect.String:
		length := uint16(len(value.(string)))
		msg.Push(length)
		pushSize = uint32(copy(tmpBuffer, value.(string)))
		break
	case reflect.Struct:
		structBuffer := bytes.Buffer{}
		err := binary.Write(&structBuffer, msg.order, value)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(structBuffer)

		pushSize = uint32(copy(tmpBuffer, structBuffer.Bytes()))

		break
	case reflect.Slice:
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

func (msg *Message) Peek(outValue interface{}) uint32{
	var peekSize uint32
	//var tmpBuffer []byte


	switch reflect.TypeOf(outValue).Kind(){
	case reflect.Ptr:
		buf := bytes.NewReader(msg.GetPayloadBuffer())
		err := binary.Read(buf, msg.order, outValue)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
		}

		fmt.Println(outValue)
		break
	case reflect.Struct:
		fmt.Println("fff")
		break
	}

	/*
	switch outValue.(type){
	case *bool, *byte:
		pOutValue := outValue.(*byte)
		*pOutValue = msg.buffer[msg.front]
		peekSize = 1
		break
	case *uint16, *int16:
		tmpBuffer = msg.buffer[msg.front : msg.front+ 2]
		pOutValue := outValue.(*uint16)
		*pOutValue = msg.order.Uint16(tmpBuffer)
		peekSize = 2
		break
	case *uint32, *int32:
		tmpBuffer = msg.buffer[msg.front : msg.front+ 4]
		pOutValue := outValue.(*uint32)
		*pOutValue = msg.order.Uint32(tmpBuffer)
		peekSize = 4
		break
	case *uint64, *int64:
		tmpBuffer = msg.buffer[msg.front : msg.front+ 8]
		pOutValue := outValue.(*uint64)
		*pOutValue = msg.order.Uint64(tmpBuffer)
		peekSize = 8
		break
	case *string:
		var length uint16
		msg.Pop(&length)
		tmpBuffer = msg.buffer[msg.front : msg.front+ uint32(length)]
		pOutValue := outValue.(*string)
		*pOutValue = string(tmpBuffer)
		peekSize = uint32(length)
		break
	case *[]byte:
		var length uint16
		msg.Pop(&length)
		tmpBuffer = msg.buffer[msg.front : msg.front+ uint32(length)]
		pOutValue := outValue.(*[]byte)
		*pOutValue = tmpBuffer
		peekSize = uint32(length)
		break
	case types.Struct:
		break
	default:
		return 0
	}
	//*/

	return peekSize
}

func (msg *Message) Pop(outValue interface{}) uint32 {
	popSize := msg.Peek(outValue)

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

func (msg *Message) EncodeRSA(clientPublicKey *rsa.PublicKey){
	if msg.isCryptoType(RSA) != true {
		//TODO_MSG :: 로그 삽입
		return
	}

	cipherMsg, err := rsa.EncryptPKCS1v15(cryptoRand.Reader, clientPublicKey, msg.buffer)
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

func (msg *Message) packHeader(header NetHeader) {
	msg.buffer[0] = byte(header.PacketType)
	msg.buffer[1] = byte(header.CryptoType)
	msg.buffer[2] = header.RandKey
	msg.order.PutUint16(msg.buffer[3:5], header.PayloadLength)
	msg.buffer[5] = header.CheckSum
}

func (msg *Message) isCryptoType(cryptoType CryptoType) bool {
	return CryptoType(msg.buffer[1]) == cryptoType
}