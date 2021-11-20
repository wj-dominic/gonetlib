package ringbuffer

import (
	"encoding/binary"
	"log"
)

const (
	BUFF_SIZE uint32 = 4096 //4kb
)

type RingBuffer struct{
	Buffer	[]byte

	Front	uint32
	Rear	uint32
	Capacity uint32

	Order	binary.ByteOrder
}

func NewRingBuffer(isLittleEndian bool, size ...uint32) *RingBuffer{
	bufferSize := BUFF_SIZE
	if len(size) > 0{
		bufferSize = size[0]
	}

	ring := RingBuffer{
		Buffer : make([]byte, bufferSize),

		Front : 0,
		Rear : 0,
		Capacity: BUFF_SIZE - 1, //한칸은 뺌

		Order : binary.LittleEndian,
	}

	if isLittleEndian == false {
		ring.Order = binary.BigEndian
	}

	return &ring
}

func (ring *RingBuffer) GetRearBuffer() []byte {
	return ring.Buffer[ring.Rear:]
}

func (ring *RingBuffer) GetFrontBuffer() []byte {
	return ring.Buffer[ring.Front : ring.Rear]
}

func (ring *RingBuffer) Write(value interface{}) uint32 {
	var pushSize uint32
	var tmpBuffer []byte
	tmpBuffer = make([]byte, BUFF_SIZE)

	switch value.(type){
	case bool, byte:
		tmpBuffer[0] = value.(byte)
		pushSize = 1
		break
	case uint16, int16:
		ring.Order.PutUint16(tmpBuffer, value.(uint16))
		pushSize = 2
		break
	case uint32, int32:
		ring.Order.PutUint32(tmpBuffer, value.(uint32))
		pushSize = 4
		break
	case uint64, int64:
		ring.Order.PutUint64(tmpBuffer, value.(uint64))
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

	emptySize := ring.getEmptySize()
	if emptySize < pushSize {
		log.Printf("no have enough space in ring buffer | emptySize[%d] pushSize[%d]", emptySize, pushSize)
		return 0
	}

	var directWriteSize uint32
	if ring.Rear < ring.Front {
		directWriteSize = uint32(len(ring.Buffer[ring.Rear:ring.Front - 1]))
	}else{
		directWriteSize = uint32(len(ring.Buffer[ring.Rear:]))
		if ring.Front == 0 {
			directWriteSize -= 1
		}
	}

	if directWriteSize < pushSize {
		copy(ring.Buffer[ring.Rear:], tmpBuffer[:directWriteSize])
		copy(ring.Buffer[:ring.Front], tmpBuffer[directWriteSize:])
	} else{
		copy(ring.Buffer[ring.Rear:], tmpBuffer)
	}

	ring.Rear = (ring.Rear + directWriteSize) % BUFF_SIZE

	return pushSize
}

func (ring *RingBuffer) Read(out_value interface{}, size uint32) uint32 {
	useSize := ring.getUseSize()
	if useSize < size {
		log.Printf("no have enough data to get in ring buffer | useSize[%d] getSize[%d]", useSize, size)
		return 0
	}

	var directReadSize uint32
	if ring.Rear < ring.Front {
		directReadSize = BUFF_SIZE - ring.Front
	} else {
		directReadSize = ring.Rear - ring.Front
	}

	firstSize := size - directReadSize
	secondSize := size - firstSize

	firstBuffer := make([]byte, firstSize)
	secondBuffer := make([]byte, secondSize)

	if directReadSize < size {
		copy(firstBuffer, ring.Buffer[ring.Front:])
		copy(secondBuffer, ring.Buffer[:ring.Rear])
	} else {
		copy(firstBuffer, ring.Buffer[ring.Front:])
	}





	var peekSize uint32
	var tmpBuffer []byte

	switch out_value.(type){
	case *bool, *byte:
		pOutValue := out_value.(*byte)
		*pOutValue = ring.Buffer[ring.Front]
		peekSize = 1
		break
	case *uint16, *int16:



		tmpBuffer = ring.Buffer[ring.Front : ring.Front + 2]
		pOutValue := out_value.(*uint16)
		*pOutValue = ring.Order.Uint16(tmpBuffer)
		peekSize = 2
		break
	case *uint32, *int32:
		tmpBuffer = ring.Buffer[ring.Front : ring.Front + 4]
		pOutValue := out_value.(*uint32)
		*pOutValue = ring.Order.Uint32(tmpBuffer)
		peekSize = 4
		break
	case *uint64, *int64:
		tmpBuffer = ring.Buffer[ring.Front : ring.Front + 8]
		pOutValue := out_value.(*uint64)
		*pOutValue = ring.Order.Uint64(tmpBuffer)
		peekSize = 8
		break
	case *string:
		tmpBuffer = ring.Buffer[ring.Front : ring.Front + uint32(size)]
		pOutValue := out_value.(*string)
		*pOutValue = string(tmpBuffer)
		peekSize = uint32(len(tmpBuffer))
		break
	case *[]byte:
		tmpBuffer = ring.Buffer[ring.Front : ring.Front + uint32(size)]
		pOutValue := out_value.(*[]byte)
		*pOutValue = tmpBuffer
		peekSize = uint32(len(tmpBuffer))
		break
	default:
		return 0
	}




}


func (ring *RingBuffer) getEmptySize() uint32{
	emptySize := uint32(0)

	if ring.Rear < ring.Front {
		emptySize = (ring.Front - ring.Rear) - 1
	} else {
		emptySize = ring.Capacity - (ring.Rear - ring.Front)
	}

	return emptySize
}

func (ring *RingBuffer) getUseSize() uint32 {
	useSize := uint32(0)

	if ring.Rear < ring.Front {
		useSize = BUFF_SIZE - (ring.Front - ring.Rear)
	} else {
		useSize = ring.Rear - ring.Front
	}

	return useSize
}