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

type Session interface {
	Start() error
	Stop() error
	Setup(uint64, net.Conn, SessionHandler, SessionEvent) error
	GetID() uint64
	Send(interface{})
	SessionMonitoringData() SessionMonitoringData
}

type SessionHandler interface {
	OnConnect(Session) error
	OnDisconnect(Session) error
	OnRecv(Session, *message.Message) error
	OnSend(Session, []byte) error
}

type SessionEvent interface {
	OnRelease(uint64, Session)
}

type releaseFlag struct {
	refCount int32
	flag     int32
}

type gonetSession struct {
	logger logger.Logger

	id   uint64
	conn net.Conn

	wg          sync.WaitGroup
	releaseFlag releaseFlag

	handler ISessionHandler
	event   ISessionEvent

	// Monitoring area
	monitoringData SessionMonitoringData
}

func newSession(logger logger.Logger) gonetSession {
	return gonetSession{
		logger: logger,

		id:   0,
		conn: nil,

		wg:          sync.WaitGroup{},
		releaseFlag: releaseFlag{0, 0},

		handler: nil,
		event:   nil,

		monitoringData: SessionMonitoringData{},
	}
}

func (session *gonetSession) Setup(id uint64, conn net.Conn, handler SessionHandler, event SessionEvent) error {
	session.id = id
	session.conn = conn
	session.handler = handler
	session.event = event
	return nil
}

func (session *gonetSession) GetID() uint64 {
	return session.id
}

func (session *gonetSession) acquire(force ...bool) bool {
	if atomic.LoadInt32(&session.releaseFlag.flag) == 1 {
		return false
	}

	if util.InterlockIncrement(&session.releaseFlag.refCount) == 1 {
		//다른 곳에서 release 중일 수 있으므로 1이면 획득 불가
		if len(force) > 0 {
			if force[0] == true {
				return true
			}
		}

		return false
	}

	return true
}

func (session *gonetSession) release() bool {
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

func (session *gonetSession) reset() {
	session.id = 0
	session.conn = nil

	//release flag 초기화는 가장 마지막에
	session.releaseFlag = releaseFlag{0, 0}
}
