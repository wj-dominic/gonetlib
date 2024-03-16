package session

import (
	"context"
	"gonetlib/logger"
	"gonetlib/message"
	"gonetlib/util"
	"net"
	"sync"
	"unsafe"
)

type ISession interface {
	Start() error
	Stop() error
	Setup(uint64, net.Conn, ISessionHandler, ISessionEvent)
	GetID() uint64
}

type ISessionHandler interface {
	OnConnect() error
	OnDisconnect() error
	OnRecv(packet *message.Message) error
	OnSend([]byte) error
}

type ISessionEvent interface {
	OnRelease(ISession)
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
	ctx         context.Context
	releaseFlag releaseFlag

	handler ISessionHandler
	event   ISessionEvent
}

func newSession(logger logger.ILogger, ctx context.Context) Session {
	return Session{
		logger: logger,

		id:   0,
		conn: nil,

		wg:          sync.WaitGroup{},
		ctx:         ctx,
		releaseFlag: releaseFlag{0, 0},

		handler: nil,
		event:   nil,
	}
}

func (session *Session) Setup(id uint64, conn net.Conn, handler ISessionHandler, event ISessionEvent) {
	session.id = id
	session.conn = conn
	session.handler = handler
	session.event = event
}

func (session *Session) GetID() uint64 {
	return session.id
}

func (session *Session) acquire() bool {
	if session.releaseFlag.flag == 1 {
		return false
	}

	util.InterlockIncrement(&session.releaseFlag.refCount)
	return true
}

func (session *Session) release() bool {
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
	session.handler = nil
	session.event = nil

	//release flag 초기화는 가장 마지막에
	session.releaseFlag = releaseFlag{0, 0}
}
