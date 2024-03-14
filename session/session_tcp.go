package session

import (
	"gonetlib/ringbuffer"
	"sync"
)

var wg sync.WaitGroup

type TcpSession struct {
	Session
	recvBuffer *ringbuffer.RingBuffer
}

func newTcpSession() ISession {
	return &TcpSession{
		recvBuffer: ringbuffer.NewRingBuffer(true),
	}
}

func (session *TcpSession) Start() bool {
	if session.conn == nil {
		return false
	}

	if session.handler != nil {
		session.handler.OnConnect()
	}

	wg.Add(1)
	go session.readAsync()

	return true
}

func (session *TcpSession) readAsync() {
	defer wg.Done()

	for {
		//빈 버퍼 획득
		buffer := session.recvBuffer.GetRearBuffer()

		recvSize, err := session.conn.Read(buffer)
		if err != nil {
			return
		}
	}
}

func (session *TcpSession) Stop() bool {
	return true
}
