package session

import (
	. "message"
	. "ringbuffer"
	"net"
)

type Session struct{
	sessionID	uint64			//세션 ID
	Conn		net.Conn		//TCP connection
	RecvBuffer	*RingBuffer		//수신 버퍼, 수신 스레드만 접근 (thread safe X)
	SendChannel chan *Message	//송신 버퍼, 송신이 필요한 모든 스레드에서 접근 (채널이어서 thread safe O)
}

func NewSession(ID uint64, conn net.Conn) *Session {
	return &Session{
		sessionID : ID,
		Conn : conn,
		RecvBuffer: NewRingBuffer(true, 300),
		SendChannel: make(chan *Message),
	}
}