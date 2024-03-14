package session

import (
	"context"
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

func CreateSessionManager(ctx context.Context, limit uint32) *SessionManager {
	return &SessionManager{
		ctx: ctx,
		pool: sync.Pool{
			New: func() interface{} {
				return newTcpSession()
			},
		},
	}
}

func (s *SessionManager) NewSession(id uint64, conn net.Conn, handler ISessionHandler) ISession {
	session := s.pool.Get().(ISession)
	if session == nil {
		return nil
	}

	session.Setup(id, conn, handler)
	return session
}
