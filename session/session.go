package session

/*
- 네트워크 영역의 클라이언트 객체
*/

import (
	. "gonetlib/message"
	. "gonetlib/ringbuffer"
	util "gonetlib/util"
	"io"
	"log"
	"net"
)

const (
	recvBufferSize uint32 = 4096
)

type Session struct{
	id	uint64			//세션 ID
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

//클라이언트 연결 시 호출하는 함수
func (session *Session) Start(sessionID uint64, connection net.Conn) bool {
	if connection == nil {
		log.Println("connection is nullptr")	//TODO : 로그
		return false
	}

	session.Reset()
	session.id = sessionID
	session.conn = connection

	go session.asyncRead()

	return true
}

func (session *Session) Close() {
	session.disconnectHandler()
}

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

//수신 스레드에서만 접근
func (session *Session) recvHandler(recvSize uint32, recvErr error) bool {
	if recvErr != nil {
		if recvErr == io.EOF {
			log.Println("connection is closed from client : ", session.conn.RemoteAddr().String())
			return false
		} else {
			log.Fatalln("read error : ", recvErr)
			return false
		}
	}

	session.recvBuffer.MoveRear(recvSize)

	return true
}

//세션 종료 함수 : 여러 스레드(accept, recv, send) 접근 가능 (스레드 세이프)
func (session *Session) disconnectHandler() {
	defer session.once.Reset()
	session.once.Do(func() {
		session.Reset()
	})
}