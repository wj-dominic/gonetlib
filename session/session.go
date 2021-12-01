package session

/*
- 네트워크 영역의 클라이언트 객체
*/

import (
	"crypto/rsa"
	. "gonetlib/message"
	. "gonetlib/netlogger"
	. "gonetlib/ringbuffer"
	util "gonetlib/util"
	"io"
	"net"
	"reflect"
)

type KeyChain struct{
	XOR		uint8
	RSA 	rsa.PublicKey
}

type Node interface {
	OnRecv(packet *Message) bool
}

type Session struct{
	id			uint64			//세션 ID
	conn		net.Conn		//TCP connection
	recvBuffer	*RingBuffer		//수신 버퍼, 수신 스레드만 접근 (thread safe X)
	sendChannel chan *Message	//송신 버퍼, 송신이 필요한 모든 스레드에서 접근 (채널이어서 thread safe O)

	keys		KeyChain

	node		Node

	once		util.Once
}


func NewSession(node Node) *Session {
	return &Session{
		id : 0,
		conn : nil,
		recvBuffer: NewRingBuffer(true, 300),
		sendChannel: make(chan *Message),

		node : node,

		keys : KeyChain{0, rsa.PublicKey{}},

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


	netHeader := NetHeader{}
	headerSize := util.Sizeof(reflect.ValueOf(netHeader))
	if headerSize == -1 {
		GetLogger().Error("header size was wrong...")
		return false
	}

	for {
		if session.recvBuffer.GetUseSize() <= uint32(headerSize) {
			break
		}

		session.recvBuffer.Peek(&netHeader, uint32(headerSize))

		packetSize := uint32(headerSize) + uint32(netHeader.PayloadLength)
		if session.recvBuffer.GetUseSize() < packetSize {
			break
		}

		session.recvBuffer.MoveFront(uint32(headerSize))

		packet := NewMessage(true)

		packet.PushHeader(&netHeader)
		session.recvBuffer.Read(packet.GetPayloadBuffer(), uint32(netHeader.PayloadLength))

		if session.onRecv(packet) == false {
			return false
		}
	}

	return true
}

//패킷 수신 이벤트 함수 : 수신 스레드에서만 접근
func (session *Session) onRecv(packet *Message) bool {
	if packet == nil {
		return false
	}

	if packet.IsValid() == false {
		//TODO 로그
		return false
	}

	cryptoType := packet.GetCryptoType()
	switch cryptoType{
	case NONE:
		break
	case XOR:
		packet.DecodeXOR(session.keys.XOR)
		break
	case RSA:
		//packet.DecodeRSA() //TODO 서버 개인 키 필요
		break
	default:
		//TODO 로그
		return false
	}

	packetType := packet.GetType()
	switch packetType{
	case SYN:
		packet.Pop(&session.keys.RSA)
		break
	case SYN_ACK:
		packet.Pop(&session.keys.XOR)
		break
	case ESTABLISHED:
		session.node.OnRecv(packet) //TODO 콘텐츠 쪽에 전달
		break
	default:
		//TODO 로그
		return false
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
