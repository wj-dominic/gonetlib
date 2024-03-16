package session

import (
	"context"
	"gonetlib/logger"
	"net"
	"sync"
)

type ISessionManager interface {
	NewSession(uint64, net.Conn, ISessionHandler) ISession
}

type SessionManager struct {
	ctx  context.Context
	pool sync.Pool
}

func CreateSessionManager(logger logger.ILogger, limit uint32, ctx context.Context) *SessionManager {
	return &SessionManager{
		ctx: ctx,
		pool: sync.Pool{
			New: func() interface{} {
				return newTcpSession(logger, ctx)
			},
		},
	}
}

//TODO: session 관리자 만들기

//TODO: session stop만들기

func (s *SessionManager) NewSession(id uint64, conn net.Conn, handler ISessionHandler) ISession {
	session := s.pool.Get().(ISession)
	if session == nil {
		return nil
	}

	session.Setup(id, conn, handler, s)
	return session
}

func (s *SessionManager) OnRelease(session ISession) {
	if session == nil {
		return
	}

	s.pool.Put(session)
}
