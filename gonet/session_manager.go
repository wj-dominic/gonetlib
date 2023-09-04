package gonet

import (
	"context"
	"gonetlib/session"
	"net"
	"sync"
)

type SessionManager struct {
	ctx  context.Context
	pool sync.Pool
}

func NewSessionManager(ctx context.Context, limit uint32) *SessionManager {
	return &SessionManager{
		ctx: ctx,
		pool: sync.Pool{
			New: func() interface{} {
				return session.NewSession()
			},
		},
	}
}

func (s *SessionManager) NewSession(context.Context, net.Conn) ISession {

}
