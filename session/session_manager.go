package session

import (
	"context"
	"net"
	"sync"
)

type ISessionManager interface {
	NewSession(net.Conn) ISession
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
				return NewSession()
			},
		},
	}
}

func (s *SessionManager) NewSession(net.Conn) ISession {
	return s.pool.Get().(ISession)
}
