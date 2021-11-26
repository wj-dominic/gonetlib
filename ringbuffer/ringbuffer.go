package ringbuffer

import (
	"encoding/binary"
	"log"
	"reflect"
)

const (
	bufferSize uint32 = 4096 //4kb
)

type RingBuffer struct{
	buffer []byte

	front    	uint32
	rear 		uint32
	cap  		uint32

	order binary.ByteOrder
}

func NewRingBuffer(isLittleEndian bool, size ...uint32) *RingBuffer{
	bufferSize := bufferSize
	if len(size) > 0{
		bufferSize = size[0]
	}

	ring := RingBuffer{
		buffer: make([]byte, bufferSize),

		front: 0,
		rear:  0,
		cap:   bufferSize - 1, //한칸은 뺌

		order: binary.LittleEndian,
	}

	if isLittleEndian == false {
		ring.order = binary.BigEndian
	}

	return &ring
}

func (ring *RingBuffer) Clear() {
	ring.front = bufferSize
	ring.rear = bufferSize
}

func (ring *RingBuffer) GetRearBuffer() []byte {
	if ring.IsFull() == true {
		return nil
	}

	directWriteSize := ring.getDirectWriteSize()

	return ring.buffer[ring.rear : ring.rear + directWriteSize]
}

func (ring *RingBuffer) GetFrontBuffer() []byte {
	if ring.IsEmpty() == true {
		return nil
	}

	directReadSize := ring.getDirectReadSize()

	return ring.buffer[ring.front : ring.front + directReadSize]
}

func (ring *RingBuffer) IsEmpty() bool {
	return ring.front == ring.rear
}

func (ring *RingBuffer) IsFull() bool {
	return ring.rear+ 1 == ring.front
}

func (ring *RingBuffer) Write(value interface{}) uint32 {
	var pushSize uint32
	var tmpBuffer []byte
	length := reflect.ValueOf(value).Len()
	tmpBuffer = make([]byte, length)

	switch value.(type){
	case bool, byte:
		tmpBuffer[0] = value.(byte)
		pushSize = 1
		break
	case uint16, int16:
		ring.order.PutUint16(tmpBuffer, value.(uint16))
		pushSize = 2
		break
	case uint32, int32:
		ring.order.PutUint32(tmpBuffer, value.(uint32))
		pushSize = 4
		break
	case uint64, int64:
		ring.order.PutUint64(tmpBuffer, value.(uint64))
		pushSize = 8
		break
	case string:
		pushSize = uint32(copy(tmpBuffer, value.(string)))
		break
	case []byte:
		pushSize = uint32(copy(tmpBuffer, value.([]byte)))
		break
	default:
		return 0
	}

	emptySize := ring.GetEmptySize()
	if emptySize < pushSize {
		log.Printf("no have enough space in ring buffer | emptySize[%d] pushSize[%d]", emptySize, pushSize)
		return 0
	}

	directWriteSize := ring.getDirectWriteSize()
	if directWriteSize < pushSize {
		copy(ring.buffer[ring.rear:], tmpBuffer[:directWriteSize])
		copy(ring.buffer[:ring.front], tmpBuffer[directWriteSize:])
	} else{
		copy(ring.buffer[ring.rear:], tmpBuffer)
	}

	ring.rear = (ring.rear + pushSize) % (ring.cap + 1)

	return pushSize
}

func (ring *RingBuffer) Read(out_value interface{}, size uint32) uint32 {
	peekSize := ring.Peek(out_value, size)

	ring.front = (ring.front + peekSize) % (ring.cap + 1)

	return peekSize
}

func (ring *RingBuffer) Peek(out_value interface{}, size uint32) uint32 {
	useSize := ring.GetUseSize()
	if useSize < size {
		log.Printf("no have enough data to get in ring buffer | useSize[%d] getSize[%d]", useSize, size)
		return 0
	}

	directReadSize := ring.getDirectReadSize()

	firstSize := directReadSize
	secondSize := uint32(0)

	if directReadSize < size {
		secondSize = size - directReadSize
	}

	firstBuffer := make([]byte, firstSize)
	secondBuffer := make([]byte, secondSize)

	copy(firstBuffer, ring.buffer[ring.front:])
	if directReadSize < size {
		copy(secondBuffer, ring.buffer[:ring.rear])
	}

	tmpBuffer := append(firstBuffer, secondBuffer...)

	var peekSize uint32

	switch out_value.(type){
	case *bool, *byte:
		pOutValue := out_value.(*byte)
		*pOutValue = ring.buffer[ring.front]
		peekSize = 1
		break
	case *uint16, *int16:
		pOutValue := out_value.(*uint16)
		*pOutValue = ring.order.Uint16(tmpBuffer)
		peekSize = 2
		break
	case *uint32, *int32:
		pOutValue := out_value.(*uint32)
		*pOutValue = ring.order.Uint32(tmpBuffer)
		peekSize = 4
		break
	case *uint64, *int64:
		pOutValue := out_value.(*uint64)
		*pOutValue = ring.order.Uint64(tmpBuffer)
		peekSize = 8
		break
	case *string:
		pOutValue := out_value.(*string)
		*pOutValue = string(tmpBuffer)
		peekSize = uint32(len(tmpBuffer))
		break
	case *[]byte:
		pOutValue := out_value.(*[]byte)
		*pOutValue = tmpBuffer
		peekSize = uint32(len(tmpBuffer))
		break
	default:
		return 0
	}

	return peekSize
}

func (ring *RingBuffer) MoveRear(offset uint32) {
	if ring.GetEmptySize() < offset {
		return
	}

	ring.rear = (ring.rear + offset) % (ring.cap + 1)
}

func (ring *RingBuffer) MoveFront(offset uint32) {
	if ring.GetUseSize() < offset {
		return
	}

	ring.front = (ring.front + offset) % (ring.cap + 1)
}

func (ring *RingBuffer) GetEmptySize() uint32{
	emptySize := uint32(0)

	if ring.rear < ring.front {
		emptySize = (ring.front - ring.rear) - 1
	} else {
		emptySize = ring.cap - (ring.rear - ring.front)
	}

	return emptySize
}

func (ring *RingBuffer) GetUseSize() uint32 {
	useSize := uint32(0)

	if ring.rear < ring.front {
		useSize = (ring.cap + 1) - (ring.front - ring.rear)
	} else {
		useSize = ring.rear - ring.front
	}

	return useSize
}

func (ring *RingBuffer) getDirectWriteSize() uint32 {
	var directWriteSize uint32

	if ring.rear < ring.front {
		directWriteSize = uint32(len(ring.buffer[ring.rear :ring.front- 1]))
	}else{
		directWriteSize = uint32(len(ring.buffer[ring.rear:]))
		if ring.front == 0 {
			directWriteSize -= 1
		}
	}

	return directWriteSize
}

func (ring *RingBuffer) getDirectReadSize() uint32 {
	var directReadSize uint32

	if ring.rear < ring.front {
		directReadSize = (ring.cap + 1) - ring.front
	} else {
		directReadSize = ring.rear - ring.front
	}

	return directReadSize
}