package session

/*
- 네트워크 영역의 클라이언트 객체
*/

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	. "gonetlib/message"
	. "gonetlib/netlogger"
	. "gonetlib/ringbuffer"
	util "gonetlib/util"
	"io"
	"net"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

const (
	maxSendBufferSize uint32 = 300
)

type keyChain struct{
	XOR		uint8
	RSA 	rsa.PublicKey
}

type Node interface {
	OnConnect()
	OnDisconnect()
	OnRecv(packet *Message) bool
	OnSend(sendBytes int)
}

type ioBlock struct {
	refCount	int32
	releaseFlag	int32
}

type Session struct{
	id			uint64			//세션 ID
	recvBuffer	*RingBuffer		//수신 버퍼, 수신 스레드만 접근 (thread safe X)
	sendChannel chan *Message	//송신 버퍼, 송신이 필요한 모든 스레드에서 접근 (채널이어서 thread safe O)
	keys 		keyChain

	socket 		net.Conn //TCP connection
	node  	 	Node

	ioblock		ioBlock

	sendOnce	util.Once
	closeOnce	util.Once
	wg 			sync.WaitGroup
}

func NewSession() *Session {
	return &Session{
		id : 0,
		recvBuffer: NewRingBuffer(true, 300),
		sendChannel: make(chan *Message, maxSendBufferSize),
		keys : keyChain{0, rsa.PublicKey{}},

		socket: nil,
		node :  nil,

		ioblock : ioBlock {0, 0},

		sendOnce: util.Once{},
		closeOnce: util.Once{},
		wg		: sync.WaitGroup{},
	}
}

func (session *Session) Setup(sessionID uint64, connection net.Conn, node Node) {
	if connection == nil {
		GetLogger().Error("connection is nullptr")
		return
	}

	session.id = sessionID
	session.socket = connection
	session.node = node
	session.ioblock.refCount = 1 	//릴리즈 방지를 위해 우선 1로 세팅
	session.ioblock.releaseFlag = 0
}

// Start : 클라이언트 연결 시 호출하는 함수 (accept 스레드에서 접근)
func (session *Session) Start() {
	session.connectHandler()
}

// Close : 클라이언트 연결 종료 함수
func (session *Session) Close() {
	if session.acquire() == false {
		return
	}

	session.closesocket()

	session.release()
}

// Reset : 세션 초기화 함수
func (session *Session) Reset() {
	session.id = 0

	session.socket = nil
	session.node = nil
	session.keys = keyChain{0, rsa.PublicKey{}}

	session.recvBuffer.Clear()
	for len(session.sendChannel) > 0 {
		select{
		case <- session.sendChannel:
			break
		default:
			break
		}
	}

	session.sendOnce.Reset()
	session.closeOnce.Reset()

	session.ioblock = ioBlock{0, 0}
}

func (session *Session) SendPost(packet *Message) bool {
	if packet == nil {
		GetLogger().Error("Failed to send | packet is nullptr")
		return false
	}

	session.sendChannel <- packet

	session.sendHandler()

	return true
}

func (session *Session) IsConnected() bool {
	return session.id != 0 && session.ioblock.releaseFlag != 1
}

//Accept 스레드에서 접근하는 함수
func (session *Session) connectHandler() {
	if session.acquire() == false {							//ref = 2
		return
	}

	defer func(){
		util.InterlockDecrement(&session.ioblock.refCount) 		//ref = 2
		session.release()                    				//ref = 1
	}()

	fmt.Println("success to connect! : ", session.id)

	session.node.OnConnect()

	if session.acquire() == true { go session.asyncRead() }	//ref = 3
}

//수신 스레드
func (session *Session) asyncRead() {
	fmt.Printf("begin async read routine... | sessionID[%d] refCount[%d] releaseFlag[%d]\n", session.id, session.ioblock.refCount, session.ioblock.releaseFlag)
	session.wg.Add(1)

	defer func() {
		session.wg.Done()
		fmt.Printf("end async read routine... | sessionID[%d] refCount[%d] releaseFlag[%d]\n", session.id, session.ioblock.refCount, session.ioblock.releaseFlag)
		session.release()
	}()

	for {
		buffer := session.recvBuffer.GetRearBuffer()

		recvSize, err := session.socket.Read(buffer)

		if session.recvHandler(uint32(recvSize), err) == false {
			break
		}
	}
}

// recvHandler : 수신 스레드에서만 접근
func (session *Session) recvHandler(recvSize uint32, recvErr error) bool {
	if recvErr != nil {
			if recvErr == io.EOF {
				GetLogger().Error("connection is closed from client : " + session.socket.RemoteAddr().String())
				return false
			} else {
				GetLogger().Error("read error : " + recvErr.Error())
			return false
		}
	}

	if session.recvBuffer.MoveRear(recvSize) == false {
		GetLogger().Error("failed to receive | recvSize[%d]", recvSize)
		return false
	}

	netHeader := NetHeader{}
	headerSize := util.Sizeof(reflect.ValueOf(netHeader))
	if headerSize == -1 {
		GetLogger().Error("header size was wrong...")
		return false
	}

	for {
		netHeader := NetHeader{}

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
		packet.MoveRear(uint32(netHeader.PayloadLength))

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
		GetLogger().Error("invalid crypto type of packet | cryptoType[%d]", cryptoType)
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
		GetLogger().Error("invalid packet type of packet | packetType[%d]", packetType)
		return false
	}

	return true
}


func (session *Session) sendHandler() {
	session.sendOnce.Do(func() {
		if session.acquire() == true { go session.asyncWrite() }
	})
}

func (session *Session) asyncWrite() {
	fmt.Printf("begin async write routine... | sessionID[%d] refCount[%d] releaseFlag[%d]\n", session.id, session.ioblock.refCount, session.ioblock.releaseFlag)
	session.wg.Add(1)

	defer func() {
		session.wg.Done()
		fmt.Printf("end async write routine... |	sessionID[%d] refCount[%d] releaseFlag[%d]\n", session.id, session.ioblock.refCount, session.ioblock.releaseFlag)
		session.release()
		session.sendOnce.Reset()
	}()

	sendBuffer := bytes.Buffer{}
	for len(session.sendChannel) > 0 {
		select{
			case msg := <- session.sendChannel:
				sendBuffer.Write(msg.GetBuffer())
				break
		}

		time.Sleep(1 * time.Millisecond)
	}

	if sendBuffer.Len() <= 0 {
		return
	}

	_ = session.socket.SetWriteDeadline(time.Now().Add(5 * time.Second))
	sendBytes, err := session.socket.Write(sendBuffer.Bytes())
	if err != nil{
		fmt.Printf("Failed to send packet to client | err[%s] sendBytes[%d]\n", err.Error(), sendBuffer.Len())
		GetLogger().Error("Failed to send packet to client | err[%s] sendBytes[%d]", err.Error(), sendBuffer.Len())
		return
	}

	session.node.OnSend(sendBytes)
}

//세션 종료 함수 : 여러 스레드(accept, recv, send) 접근 가능 (스레드 세이프하게 만들어야 함)
func (session *Session) disconnectHandler() {
	if session.canDisconnect() == false {
		return
	}

	session.wg.Wait()

	session.closesocket()

	session.node.OnDisconnect()

	session.Reset()
}

// acquire 세션 획득 메소드
func (session *Session) acquire() bool {
	refCount := util.InterlockIncrement(&session.ioblock.refCount)
	if refCount == 1 {
		//릴리즈 중인 세션이므로 릴리즈 수행
		session.release()
		return false
	}

	if util.InterlockedCompareExchange(&session.ioblock.releaseFlag, 1, 1) == true {
		//릴리즈 중인 세션이므로 릴리즈 수행
		session.release()
		return false
	}

	return true
}

// release 세션 반환 메소드
func (session *Session) release() {
	refCount := util.InterlockDecrement(&session.ioblock.refCount)
	if refCount == 0 {
		//릴리즈
		session.disconnectHandler()
		return

	} else if refCount < 0 {
		//fatal : 문제가 심각함
		GetLogger().Error("session refer count is minus | refCount[%d]", refCount)
		return
	}
}

func (session *Session) canDisconnect() bool {
	destIOBlock := (*int64)(unsafe.Pointer(&ioBlock{0, 1}))
	compareIOBlock := (*int64)(unsafe.Pointer(&ioBlock{0, 0}))

	originIOBlock := (*int64)(unsafe.Pointer(&session.ioblock))

	if util.InterlockedCompareExchange64(originIOBlock, *destIOBlock, *compareIOBlock) == false {
		GetLogger().Debug("Can't release | originBlock[%d]", *originIOBlock)
		fmt.Printf("Can't release | originBlock[%d]\n", *originIOBlock)
		return false
	}

	return true
}

func (session *Session) closesocket() {
	session.closeOnce.Do(func(){
		session.socket.Close()
		close(session.sendChannel)
	})
}
