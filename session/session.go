package session

import (
	"context"
	"gonetlib/logger"
	"gonetlib/message"
	"net"
	"sync"
)

type ISession interface {
	Start() error
	Stop() error
	Setup(uint64, net.Conn, ISessionHandler)
	GetID() uint64
}

type ISessionHandler interface {
	OnConnect() error
	OnDisconnect() error
	OnRecv(packet *message.Message) error
	OnSend([]byte) error
}

type Session struct {
	id       uint64
	conn     net.Conn
	handler  ISessionHandler
	refCount int32
	wg       sync.WaitGroup
	ctx      context.Context
	logger   logger.ILogger
}

func newSession(logger logger.ILogger, ctx context.Context) Session {
	return Session{
		id:       0,
		conn:     nil,
		handler:  nil,
		refCount: 0,
		wg:       sync.WaitGroup{},
		ctx:      ctx,
		logger:   logger,
	}
}

func (session *Session) Setup(id uint64, conn net.Conn, handler ISessionHandler) {
	session.id = id
	session.conn = conn
	session.handler = handler
}

func (session *Session) GetID() uint64 {
	return session.id
}
