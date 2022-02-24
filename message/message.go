package message

import (
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"fmt"
	. "gonetlib/netlogger"
	"gonetlib/util"
	"log"
	mathRand "math/rand"
	"reflect"
)

type PacketType uint8

const (
	SYN         PacketType = 1 + iota //공개키 주고 받음, 암호화 없음
	SYN_ACK                           //패킷 코드, 키 주고 받음 (XOR 용), 이 패킷은 RSA 암호화가 기본
	ESTABLISHED                       //이 패킷은 연결된 후 패킷, 여기서부터는 서로 공개키와 패킷 코드, 키를 알고 있으므로 인코딩은 선택하면됨
)

type CryptoType uint8

const (
	NONE CryptoType = 0 + iota //암호화 없음
	XOR                        //빠른 암호화
	RSA                        //느린 암호화
)

type NetHeader struct {
	PacketType    PacketType
	CryptoType    CryptoType
	RandKey       uint8
	PayloadLength uint16
	CheckSum      uint8
}

const (
	headerSize  uint32 = 6
	payloadSize uint32 = 300
	bufferSize  uint32 = headerSize + payloadSize
)

type Message struct {
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

func (msg *Message) GetBuffer() []byte {
	return msg.buffer[:msg.rear]
}

func (msg *Message) GetHeaderBuffer() []byte {
	return msg.buffer[:headerSize]
}

func (msg *Message) GetPayloadBuffer() []byte {
	return msg.buffer[msg.front:]
}

func (msg *Message) GetPayloadLength() uint32 {
	return msg.rear - msg.front
}

func (msg *Message) GetLength() uint32 {
	return headerSize + msg.GetPayloadLength()
}

func (msg *Message) SetHeader(packetType PacketType, cryptoType CryptoType) {
	var header NetHeader
	header.PacketType = packetType
	header.CryptoType = cryptoType
	header.RandKey = uint8(mathRand.Intn(256))
	header.PayloadLength = uint16(msg.GetPayloadLength())
	header.CheckSum = msg.generateChecksum()

	msg.packHeader(header) //pragma pack(1)
}

func (msg *Message) PushHeader(header *NetHeader) {
	msg.packHeader(*header)
}

func (msg *Message) Push(value interface{}) uint32 {
	pushSize := uint32(util.Sizeof(reflect.ValueOf(value)))

	if msg.getFreeLength() < pushSize {
		fmt.Println(value, pushSize)
		return 0
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Bool:
		boolValue := value.(bool)
		if boolValue == true {
			msg.buffer[msg.rear] = 1
		} else {
			msg.buffer[msg.rear] = 0
		}
		break
	case reflect.Int8:
		tempValue := value.(int8)
		msg.buffer[msg.rear] = uint8(tempValue)
		break
	case reflect.Uint8:
		msg.buffer[msg.rear] = value.(byte)
		break
	case reflect.Int16:
		tempValue := value.(int16)
		msg.order.PutUint16(msg.buffer[msg.rear:], uint16(tempValue))
		break
	case reflect.Uint16:
		msg.order.PutUint16(msg.buffer[msg.rear:], value.(uint16))
		break
	case reflect.Int32:
		tempValue := value.(int32)
		msg.order.PutUint32(msg.buffer[msg.rear:], uint32(tempValue))
		break
	case reflect.Uint32:
		msg.order.PutUint32(msg.buffer[msg.rear:], value.(uint32))
		break
	case reflect.Int64:
		tempValue := value.(int64)
		msg.order.PutUint64(msg.buffer[msg.rear:], uint64(tempValue))
		break
	case reflect.Uint64:
		msg.order.PutUint64(msg.buffer[msg.rear:], value.(uint64))
		break
	case reflect.Int:
		tempValue := value.(int)
		if pushSize == 8 {
			msg.order.PutUint64(msg.buffer[msg.rear:], uint64(tempValue))
		} else {
			msg.order.PutUint32(msg.buffer[msg.rear:], uint32(tempValue))
		}
		break
	case reflect.Uint:
		tempValue := value.(uint)
		if pushSize == 8 {
			msg.order.PutUint64(msg.buffer[msg.rear:], uint64(tempValue))
		} else {
			msg.order.PutUint32(msg.buffer[msg.rear:], uint32(tempValue))
		}
		break
	case reflect.String:
		length := uint16(len(value.(string)))
		msg.Push(length)
		copy(msg.buffer[msg.rear:], value.(string))
		pushSize = uint32(length)
		break
	case reflect.Struct:
		target := reflect.ValueOf(value)
		for i, n := 0, target.NumField() ; i < n ; i++ {
			msg.Push(target.Field(i).Interface())
		}
		pushSize = 0	//이미 위에서 넣기 때문에 pushsize는 0으로 변경
		break
	case reflect.Slice:
		length := uint16(len(value.([]byte)))
		msg.Push(length)
		copy(msg.buffer[msg.rear:], value.([]byte))
		break
	default:
		return 0
	}

	msg.rear += pushSize

	return pushSize
}

func (msg *Message) Peek(outValue interface{}) uint32 {
	peekSize := uint32(util.Sizeof(reflect.ValueOf(outValue).Elem()))

	switch reflect.TypeOf(outValue).Kind() {
	case reflect.Ptr:
		switch reflect.TypeOf(outValue).Elem().Kind() {
		case reflect.Bool:
			pOutValue := outValue.(*bool)
			tempValue := msg.buffer[msg.front]
			if tempValue == 1 {
				*pOutValue = true
			} else if tempValue == 0 {
				*pOutValue = false
			} else {
				GetLogger().Error("peeked value is not boolean : " + string(tempValue))
			}
			break

		case reflect.Int8:
			pOutValue := outValue.(*int8)
			*pOutValue = int8(msg.buffer[msg.front])
			break

		case reflect.Uint8:
			pOutValue := outValue.(*uint8)
			*pOutValue = msg.buffer[msg.front]
			break

		case reflect.Int16:
			pOutValue := outValue.(*int16)
			*pOutValue = int16(msg.order.Uint16(msg.GetPayloadBuffer()))
			break

		case reflect.Uint16:
			pOutValue := outValue.(*uint16)
			*pOutValue = msg.order.Uint16(msg.GetPayloadBuffer())
			break

		case reflect.Int32:
			pOutValue := outValue.(*int32)
			*pOutValue = int32(msg.order.Uint32(msg.GetPayloadBuffer()))
			break

		case reflect.Uint32:
			pOutValue := outValue.(*uint32)
			*pOutValue = msg.order.Uint32(msg.GetPayloadBuffer())
			break

		case reflect.Int64:
			pOutValue := outValue.(*int64)
			*pOutValue = int64(msg.order.Uint64(msg.GetPayloadBuffer()))
			break

		case reflect.Uint64:
			pOutValue := outValue.(*uint64)
			*pOutValue = msg.order.Uint64(msg.GetPayloadBuffer())
			break

		case reflect.Int:
			pOutValue := outValue.(*int)

			if peekSize == 8 {
				*pOutValue = int(msg.order.Uint64(msg.GetPayloadBuffer()))
			} else {
				*pOutValue = int(msg.order.Uint32(msg.GetPayloadBuffer()))
			}

			break
		case reflect.Uint:
			pOutValue := outValue.(*uint)

			if peekSize == 8 {
				*pOutValue = uint(msg.order.Uint64(msg.GetPayloadBuffer()))
			} else {
				*pOutValue = uint(msg.order.Uint32(msg.GetPayloadBuffer()))
			}
			break

		case reflect.Struct:
			target := reflect.ValueOf(outValue).Elem()
			tempPeekSize := uint32(0)
			for i, n := 0, target.NumField() ; i < n ; i++ {
				fieldPeekSize := msg.Peek(target.Field(i).Addr().Interface())
				msg.front += fieldPeekSize
				tempPeekSize += fieldPeekSize
			}

			if tempPeekSize != peekSize {
				fmt.Printf("defer peek size | tempPeekSize[%d] peekSize[%d]", tempPeekSize, peekSize)
				return 0
			}

			msg.front -= tempPeekSize
			break

		case reflect.String:
			var length uint16
			tempPeekSize := msg.Peek(&length)
			tmpBuffer := msg.buffer[msg.front + tempPeekSize : msg.front+ tempPeekSize + uint32(length)]
			pOutValue := outValue.(*string)
			*pOutValue = string(tmpBuffer)
			peekSize = uint32(length) + tempPeekSize
			break

		case reflect.Slice:
			var length uint16
			msg.Pop(&length)
			pOutValue := outValue.(*[]byte)
			*pOutValue = msg.buffer[msg.front : msg.front+uint32(length)]
			peekSize = uint32(length)
			break
		}
		break
	}

	return peekSize
}

func (msg *Message) Pop(outValue interface{}) uint32 {
	popSize := msg.Peek(outValue)

	msg.front += popSize

	return popSize
}

func (msg *Message) EncodeXOR(key uint8) {
	if msg.isCryptoType(XOR) != true {
		//TODO_MSG :: 로그 삽입
		return
	}

	randKey := msg.buffer[2]
	dstBuffer := msg.buffer[headerSize-1 : msg.rear]

	num := uint32(1)
	for i := range dstBuffer {
		p := dstBuffer[i] ^ uint8(uint32(randKey)+num)
		dstBuffer[i] = p ^ uint8(uint32(key)+num)
	}
}

func (msg *Message) DecodeXOR(key uint8) {
	if msg.isCryptoType(XOR) != true {
		//TODO_MSG :: 로그 삽입
		return
	}

	randKey := msg.buffer[2]
	dstBuffer := msg.buffer[headerSize-1 : msg.rear]

	num := uint32(1)
	for i := range dstBuffer {
		p := dstBuffer[i] ^ uint8(uint32(key)+num)
		dstBuffer[i] = p ^ uint8(uint32(randKey)+num)
	}

	//체크섬 확인
	recvChecksum := dstBuffer[0]
	generatedChecksum := msg.generateChecksum()
	if recvChecksum != generatedChecksum {
		//TODO_MSG :: 로그 삽입
		log.Fatalln("mismatch check sum : ", recvChecksum, generatedChecksum)
		return
	}
}

func (msg *Message) EncodeRSA(clientPublicKey *rsa.PublicKey) {
	if msg.isCryptoType(RSA) != true {
		//TODO_MSG :: 로그 삽입
		return
	}

	cipherMsg, err := rsa.EncryptPKCS1v15(cryptoRand.Reader, clientPublicKey, msg.buffer)
	if err != nil {
		//TODO_MSG :: 로그 추가 필요
		return
	}

	cipherMsgLength := len(cipherMsg)

	msg.clear()
	copy(msg.buffer, cipherMsg)
	msg.rear += uint32(cipherMsgLength)
}

func (msg *Message) DecodeRSA(servPrivateKey *rsa.PrivateKey) {
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

func (msg *Message) MoveFront(offset uint32) {
	msg.front += offset
}

func (msg *Message) MoveRear(offset uint32) {
	msg.rear += offset
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

func (msg *Message) generateChecksum() uint8 {
	var total uint32

	payload := msg.buffer[msg.front:msg.rear]

	for i := range payload {
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

func (msg *Message) IsValid() bool {
	packetType := PacketType(msg.buffer[0])
	cryptoType := CryptoType(msg.buffer[1])

	switch packetType {
	case SYN:
		if cryptoType != NONE {
			GetLogger().Error("packet is invalid | cryptoType[%d] packetType[%d]", cryptoType, packetType)
			return false
		}
		break
	case SYN_ACK:
		if cryptoType != RSA {
			GetLogger().Error("packet is invalid | cryptoType[%d] packetType[%d]", cryptoType, packetType)
			return false
		}
		break
	case ESTABLISHED:
		break
	}
	return true
}

func (msg *Message) GetType() PacketType {
	return PacketType(msg.buffer[0])
}

func (msg *Message) GetCryptoType() CryptoType {
	return CryptoType(msg.buffer[1])
}
