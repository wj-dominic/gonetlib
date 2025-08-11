package message

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"reflect"

	"github.com/wj-dominic/gonetlib/util"
)

type IMessageEncoder interface {
	Encode(key interface{}, buf []byte) error
}

type IMessageDecoder interface {
	Decode(key interface{}, buf []byte) error
}

type Message struct {
	buffer []byte

	front uint16
	rear  uint16

	order binary.ByteOrder

	encoder IMessageEncoder
	decoder IMessageDecoder
}

const (
	HeaderSize  uint16 = 5
	PayloadSize uint16 = 1020
	BufferSize  uint16 = HeaderSize + PayloadSize
)

// packet type
const (
	_ uint8 = iota
	Default
	SetKey
)

func NewMessage() *Message {
	return &Message{
		buffer: make([]byte, BufferSize),

		front: HeaderSize,
		rear:  HeaderSize,

		order: binary.LittleEndian,

		encoder: nil,
		decoder: nil,
	}
}

func (msg *Message) LittleEndian() *Message {
	msg.order = binary.LittleEndian
	return msg
}

func (msg *Message) BigEndian() *Message {
	msg.order = binary.BigEndian
	return msg
}

func (msg *Message) Encoder(encoder IMessageEncoder) *Message {
	msg.encoder = encoder
	return msg
}

func (msg *Message) Decoder(decoder IMessageDecoder) *Message {
	msg.decoder = decoder
	return msg
}

func (msg *Message) GetBuffer() []byte {
	return msg.buffer[:msg.GetSize()]
}

func (msg *Message) GetSize() uint16 {
	return msg.rear
}

func (msg *Message) GetHeaderBuffer() []byte {
	return msg.buffer[:msg.GetHeaderSize()]
}

func (msg *Message) GetHeaderSize() uint16 {
	return msg.front
}

func (msg *Message) GetPayloadBuffer() []byte {
	return msg.buffer[msg.GetHeaderSize():]
}

func (msg *Message) GetPayloadSize() uint16 {
	return msg.GetSize() - msg.GetHeaderSize()
}

func (msg *Message) setHeader(Type uint8) {
	msg.buffer[0] = Type
	msg.order.PutUint16(msg.buffer[1:3], uint16(msg.GetPayloadSize()))
	msg.buffer[3] = uint8(rand.Intn(256))
	msg.buffer[4] = msg.generateChecksum()
}

func (msg *Message) MakeHeader() {
	msg.setHeader(Default)
}

func (msg *Message) GetType() uint8 {
	return msg.buffer[0]
}

func (msg *Message) GetExpectedPayloadSize() uint16 {
	return msg.order.Uint16(msg.buffer[1:3])
}

func (msg *Message) GetChecksum() uint8 {
	return msg.buffer[4]
}

func (msg *Message) Push(value interface{}) uint16 {
	pushSize := uint16(util.Sizeof(reflect.ValueOf(value)))

	if msg.getFreeLength() < pushSize {
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
		pushSize = length
		break
	case reflect.Struct:
		target := reflect.ValueOf(value)
		for i, n := 0, target.NumField(); i < n; i++ {
			msg.Push(target.Field(i).Interface())
		}
		pushSize = 0 //이미 위에서 넣기 때문에 pushsize는 0으로 변경
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

func (msg *Message) Peek(outValue interface{}) uint16 {
	peekSize := uint16(util.Sizeof(reflect.ValueOf(outValue).Elem()))

	switch reflect.TypeOf(outValue).Kind() {
	case reflect.Ptr:
		switch reflect.TypeOf(outValue).Elem().Kind() {
		case reflect.Bool:
			pOutValue := outValue.(*bool)
			tempValue := msg.buffer[msg.front]
			if tempValue == 1 {
				*pOutValue = true
			} else {
				*pOutValue = false
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
			tempPeekSize := uint16(0)
			for i, n := 0, target.NumField(); i < n; i++ {
				fieldPeekSize := msg.Peek(target.Field(i).Addr().Interface())
				msg.front += fieldPeekSize
				tempPeekSize += fieldPeekSize
			}

			msg.front -= tempPeekSize
			peekSize = tempPeekSize
			break

		case reflect.String:
			var length uint16
			tempPeekSize := msg.Peek(&length)
			tmpBuffer := msg.buffer[msg.front+tempPeekSize : msg.front+tempPeekSize+length]
			pOutValue := outValue.(*string)
			*pOutValue = string(tmpBuffer)
			peekSize = length + tempPeekSize
			break

		case reflect.Slice:
			var length uint16
			msg.Pop(&length)
			pOutValue := outValue.(*[]byte)
			*pOutValue = msg.buffer[msg.front : msg.front+length]
			peekSize = length
			break
		}
		break
	}

	return peekSize
}

func (msg *Message) Pop(outValue interface{}) uint16 {
	popSize := msg.Peek(outValue)

	msg.front += popSize

	return popSize
}

func (msg *Message) Encode(key uint32) error {
	return msg.encoder.Encode(key, msg.buffer[3:msg.GetSize()])
}

func (msg *Message) Decode(key uint32) error {
	if err := msg.decoder.Decode(key, msg.buffer[3:msg.GetSize()]); err != nil {
		return err
	}

	recvChecksum := msg.GetChecksum()
	generatedChecksum := msg.generateChecksum()
	if recvChecksum != generatedChecksum {
		return fmt.Errorf("invalid check sum, mismatch | recv[%d] gen[%d]", recvChecksum, generatedChecksum)
	}

	return nil
}

func (msg *Message) MoveFront(offset uint16) {
	msg.front += offset
}

func (msg *Message) MoveRear(offset uint16) {
	msg.rear += offset
}

func (msg *Message) Reset() {
	msg.front = HeaderSize
	msg.rear = HeaderSize
}

func (msg *Message) getFreeLength() uint16 {
	tempFront := msg.front
	tempRear := msg.rear

	return (PayloadSize - 1) - (tempRear - tempFront)
}

func (msg *Message) generateChecksum() uint8 {
	var total uint32

	payload := msg.GetPayloadBuffer()
	for i := 0; i < int(msg.GetPayloadSize()); i++ {
		total += uint32(payload[i])
	}

	return uint8(total % 256)
}
