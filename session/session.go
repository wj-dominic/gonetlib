package session

/*
- 네트워크 영역의 클라이언트 객체
*/

import (
	. "gonetlib/logger"
	. "gonetlib/message"
	. "gonetlib/ringbuffer"
	util "gonetlib/util"
	"io"
	"net"
	"reflect"
)

type Session struct{
	id			uint64			//세션 ID
	conn		net.Conn		//TCP connection
	recvBuffer	*RingBuffer		//수신 버퍼, 수신 스레드만 접근 (thread safe X)
	sendChannel chan *Message	//송신 버퍼, 송신이 필요한 모든 스레드에서 접근 (채널이어서 thread safe O)

	once		util.Once
}

func NewSession() *Session {
	return &Session{
		id : 0,
		conn : nil,
		recvBuffer: NewRingBuffer(true, 300),
		sendChannel: make(chan *Message),
	}
}

// Start : 클라이언트 연결 시 호출하는 함수
func (session *Session) Start(sessionID uint64, connection net.Conn) bool {
	if connection == nil {
		GetLogger().Error("connection is nullptr")
		return false
	}

	session.Reset()
	session.id = sessionID
	session.conn = connection

	go session.asyncRead()

	return true
}

// Close : 클라이언트 연결 종료 함수
func (session *Session) Close() {
	session.disconnectHandler()
}

// Reset : 세션 초기화 함수
func (session *Session) Reset() {
	session.id = 0
	session.conn = nil
	session.recvBuffer.Clear()

	for len(session.sendChannel) > 0 {
		select{
		case <- session.sendChannel:
			break
		default:
			break
		}
	}

}

//수신 스레드
func (session *Session) asyncRead() {
	if session.conn == nil {
		return
	}

	defer session.disconnectHandler()

	for {
		buffer := session.recvBuffer.GetRearBuffer()

		recvSize, err := session.conn.Read(buffer)

		if session.recvHandler(uint32(recvSize), err) == false {
			break
		}
	}
}

// recvHandler : 수신 스레드에서만 접근
func (session *Session) recvHandler(recvSize uint32, recvErr error) bool {
	if recvErr != nil {
		if recvErr == io.EOF {
			GetLogger().Error("connection is closed from client : " + session.conn.RemoteAddr().String())
			return false
		} else {
			GetLogger().Error("read error : " + recvErr.Error())
			return false
		}
	}

	if session.recvBuffer.MoveRear(recvSize) == false {
		GetLogger().Error("failed to receive : " + string(recvSize))
		return false
	}


	for {
		var netHeader NetHeader
		headerSize := util.Sizeof(reflect.ValueOf(netHeader))
		if headerSize == -1 {
			GetLogger().Error("header size was wrong...")
			return false
		}

		if session.recvBuffer.GetUseSize() <= uint32(headerSize) {
			break
		}

		session.recvBuffer.Peek(&netHeader, uint32(headerSize))


	}



	return true
}

//세션 종료 함수 : 여러 스레드(accept, recv, send) 접근 가능 (스레드 세이프)
func (session *Session) disconnectHandler() {
	defer session.once.Reset()

	session.once.Do(func() {
		session.Reset()
	})
}