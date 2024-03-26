package session

import (
	"gonetlib/logger"
	"gonetlib/message"
	"gonetlib/util"
	"net"
	"sync"
	"sync/atomic"
	"unsafe"
)

type ISession interface {
	Start() error
	Stop() error
	Setup(uint64, net.Conn, ISessionHandler, ISessionEvent) error
	GetID() uint64
	Send(interface{})
}

type ISessionHandler interface {
	OnConnect(ISession) error
	OnDisconnect(ISession) error
	OnRecv(ISession, *message.Message) error
	OnSend(ISession, []byte) error
}

type ISessionEvent interface {
	OnRelease(uint64, ISession)
}

type releaseFlag struct {
	refCount int32
	flag     int32
}

type Session struct {
	logger logger.ILogger

	id   uint64
	conn net.Conn

	wg          sync.WaitGroup
	releaseFlag releaseFlag

	handler ISessionHandler
	event   ISessionEvent
}

func newSession(logger logger.ILogger) Session {
	return Session{
		logger: logger,

		id:   0,
		conn: nil,

		wg:          sync.WaitGroup{},
		releaseFlag: releaseFlag{0, 0},

		handler: nil,
		event:   nil,
	}
}

func (session *Session) Setup(id uint64, conn net.Conn, handler ISessionHandler, event ISessionEvent) error {
	session.id = id
	session.conn = conn
	session.handler = handler
	session.event = event
	return nil
}

func (session *Session) GetID() uint64 {
	return session.id
}

func (session *Session) acquire(force ...bool) bool {
	if atomic.LoadInt32(&session.releaseFlag.flag) == 1 {
		return false
	}

	if util.InterlockIncrement(&session.releaseFlag.refCount) == 1 {
		//다른 곳에서 release 중일 수 있으므로 1이면 획득 불가
		if force[0] == true {
			return true
		}

		return false
	}

	return true
}

func (session *Session) release() bool {
	if atomic.LoadInt32(&session.releaseFlag.flag) == 1 {
		return false
	}

	refCount := util.InterlockDecrement(&session.releaseFlag.refCount)
	if refCount == 0 {
		exchange := (*int64)(unsafe.Pointer(&releaseFlag{0, 1}))
		compare := (*int64)(unsafe.Pointer(&releaseFlag{0, 0}))
		origin := (*int64)(unsafe.Pointer(&session.releaseFlag))

		if util.InterlockedCompareExchange64(origin, *exchange, *compare) == true {
			return true
		}
	}

	return false
}

func (session *Session) reset() {
	session.id = 0
	session.conn = nil

	//release flag 초기화는 가장 마지막에
	session.releaseFlag = releaseFlag{0, 0}
}
