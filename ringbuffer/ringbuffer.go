package ringbuffer

const (
	MAX_SIZE uint32 = 4096 //4kb
)

type RingBuffer struct{
	Buffer	[]byte

	Front	uint32
	Rear	uint32
}

func NewRingBuffer(size ...uint32) *RingBuffer{
	bufferSize := MAX_SIZE
	if len(size) > 0{
		bufferSize = size[0]
	}

	return &RingBuffer{
		Buffer : make([]byte, bufferSize),

		Front : 0,
		Rear : 0,
	}
}

func (ring *RingBuffer) GetRearBuffer() []byte {
	return ring.Buffer[ring.Rear:]
}

func (ring *RingBuffer) GetFrontBuffer() []byte {
	return ring.Buffer[ring.Front : ring.Rear]
}

