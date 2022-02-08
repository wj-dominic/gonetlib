package ringbuffer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	. "gonetlib/netlogger"
	"gonetlib/util"
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
	ring.front = 0
	ring.rear = 0
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
	if reflect.TypeOf(value).Kind() == reflect.String {
		pushSize = uint32(len(value.(string)))
	} else if reflect.TypeOf(value).Kind() == reflect.Slice{
		pushSize = uint32(len(value.([]byte)))
	} else {
		pushSize = uint32(util.Sizeof(reflect.ValueOf(value)))
	}

	emptySize := ring.GetEmptySize()
	if emptySize < pushSize {
		GetLogger().Error("no have enough space in ring buffer")
		return 0
	}

	tmpBuffer := make([]byte, pushSize)

	switch reflect.TypeOf(value).Kind() {
	case reflect.Bool:
		boolValue := value.(bool)
		if boolValue == true {
			tmpBuffer[0] = 1
		} else {
			tmpBuffer[0] = 0
		}
		break
	case reflect.Int8:
		tempValue := value.(int8)
		tmpBuffer[0] = uint8(tempValue)
		break
	case reflect.Uint8:
		tmpBuffer[0] = value.(byte)
		break
	case reflect.Int16:
		tempValue := value.(int16)
		ring.order.PutUint16(tmpBuffer, uint16(tempValue))
		break
	case reflect.Uint16:
		ring.order.PutUint16(tmpBuffer, value.(uint16))
		break
	case reflect.Int32:
		tempValue := value.(int32)
		ring.order.PutUint32(tmpBuffer, uint32(tempValue))
		break
	case reflect.Uint32:
		ring.order.PutUint32(tmpBuffer, value.(uint32))
		break
	case reflect.Int64:
		tempValue := value.(int64)
		ring.order.PutUint64(tmpBuffer, uint64(tempValue))
		break
	case reflect.Uint64:
		ring.order.PutUint64(tmpBuffer, value.(uint64))
		break
	case reflect.Int:
		tempValue := value.(int)
		if pushSize == 8 {
			ring.order.PutUint64(tmpBuffer, uint64(tempValue))
		} else {
			ring.order.PutUint32(tmpBuffer, uint32(tempValue))
		}
		break
	case reflect.Uint:
		tempValue := value.(uint)
		if pushSize == 8 {
			ring.order.PutUint64(tmpBuffer, uint64(tempValue))
		} else {
			ring.order.PutUint32(tmpBuffer, uint32(tempValue))
		}
		break
	case reflect.String:
		copy(tmpBuffer, value.(string))
		break
	case reflect.Struct:
		structBuffer := bytes.Buffer{}
		err := binary.Write(&structBuffer, ring.order, value)
		if err != nil {
			log.Fatalln(err)
		}
		tmpBuffer = nil
		tmpBuffer = structBuffer.Bytes()
		break
	case reflect.Slice:
		tmpBuffer = nil
		tmpBuffer = value.([]byte)
		break
	default:
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

func (ring *RingBuffer) Read(outValue interface{}, size uint32) uint32 {
	peekSize := ring.Peek(outValue, size)

	ring.front = (ring.front + peekSize) % (ring.cap + 1)

	return peekSize
}

func (ring *RingBuffer) Peek(outValue interface{}, size uint32) uint32 {
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
	if reflect.TypeOf(outValue).Elem().Kind() != reflect.String && reflect.TypeOf(outValue).Elem().Kind() != reflect.Slice{
		peekSize = uint32(util.Sizeof(reflect.ValueOf(outValue)))
	}

	switch reflect.TypeOf(outValue).Kind(){
	case reflect.Ptr:
		switch reflect.TypeOf(outValue).Elem().Kind() {
		case reflect.Bool:
			pOutValue := outValue.(*bool)
			tempValue := ring.buffer[ring.front]
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
			*pOutValue = int8(ring.buffer[ring.front])
			break

		case reflect.Uint8:
			pOutValue := outValue.(*uint8)
			*pOutValue = ring.buffer[ring.front]
			break

		case reflect.Int16:
			pOutValue := outValue.(*int16)
			*pOutValue = int16(ring.order.Uint16(tmpBuffer))
			break

		case reflect.Uint16:
			pOutValue := outValue.(*uint16)
			*pOutValue = ring.order.Uint16(tmpBuffer)
			break

		case reflect.Int32:
			pOutValue := outValue.(*int32)
			*pOutValue = int32(ring.order.Uint32(tmpBuffer))
			break

		case reflect.Uint32:
			pOutValue := outValue.(*uint32)
			*pOutValue = ring.order.Uint32(tmpBuffer)
			break

		case reflect.Int64:
			pOutValue := outValue.(*int64)
			*pOutValue = int64(ring.order.Uint64(tmpBuffer))
			break

		case reflect.Uint64:
			pOutValue := outValue.(*uint64)
			*pOutValue = ring.order.Uint64(tmpBuffer)
			break

		case reflect.Int:
			pOutValue := outValue.(*int)

			if peekSize == 8 {
				*pOutValue = int(ring.order.Uint64(tmpBuffer))
			} else {
				*pOutValue = int(ring.order.Uint32(tmpBuffer))
			}

			break
		case reflect.Uint:
			pOutValue := outValue.(*uint)

			if peekSize == 8 {
				*pOutValue = uint(ring.order.Uint64(tmpBuffer))
			} else {
				*pOutValue = uint(ring.order.Uint32(tmpBuffer))
			}
			break

		case reflect.Struct:
			buf := bytes.NewReader(tmpBuffer)
			err := binary.Read(buf, ring.order, outValue)
			if err != nil {
				fmt.Println("binary.Read failed:", err)
			}
			peekSize = uint32(util.Sizeof(reflect.ValueOf(outValue).Elem()))
			break

		case reflect.String:
			pOutValue := outValue.(*string)
			*pOutValue = string(tmpBuffer)
			peekSize = size
			break

		case reflect.Slice:
			pOutValue := outValue.(*[]byte)
			*pOutValue = tmpBuffer
			peekSize = size
			break
		}
		break
	case reflect.Slice:
		peekSize = uint32(copy(outValue.([]byte), tmpBuffer[:size]))
		break
	}

	return peekSize
}

func (ring *RingBuffer) MoveRear(offset uint32) bool {
	if ring.GetEmptySize() < offset {
		return false
	}

	ring.rear = (ring.rear + offset) % (ring.cap + 1)

	return true
}

func (ring *RingBuffer) MoveFront(offset uint32) bool {
	if ring.GetUseSize() < offset {
		return false
	}

	ring.front = (ring.front + offset) % (ring.cap + 1)

	return true
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
		directWriteSize = (ring.front - 1) - ring.rear
	}else{
		directWriteSize = ring.cap - ring.rear
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