package netserver

import (
	"encoding/binary"
)

type bytes	[]byte

const (
	HEADER_SIZE		= 5
	PAYLOAD_SIZE	= 300
	MAX_SIZE		= HEADER_SIZE + PAYLOAD_SIZE
)

type Message struct{
	Buffer	bytes

	Front	uint32
	Rear	uint32

	Order	binary.ByteOrder
}

func NewMessage(isLittleEndian bool) *Message {
	var msg = Message{
		Buffer: make(bytes, MAX_SIZE),

		Front: HEADER_SIZE,
		Rear:  HEADER_SIZE,

		Order: binary.LittleEndian,
	}
	if isLittleEndian == false {
		msg.Order = binary.BigEndian
	}

	return &msg
}

func (msg *Message) GetData() bytes{
	return msg.Buffer
}

func (msg *Message) GetLength() int {
	return len(msg.Buffer)
}

func (msg *Message) Push(value interface{}) uint32 {
	var pushSize uint32
	var tmpBuffer bytes
	tmpBuffer = make(bytes, MAX_SIZE)

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
	case bytes:
		length := uint16(len(value.(bytes)))
		msg.Push(length)
		pushSize = uint32(copy(tmpBuffer, value.(bytes)))
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
	var tmpBuffer bytes

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
	case *bytes:
		var length uint16
		msg.Pop(&length)
		tmpBuffer = msg.Buffer[msg.Front : msg.Front + uint32(length)]
		pOutValue := out_value.(*bytes)
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


func (msg *Message) Encode() {

}


func (msg *Message) _GetFreeLength() uint32 {
	tempFront := msg.Front
	tempRear := msg.Rear

	return (PAYLOAD_SIZE - 1) - (tempRear - tempFront)
}
