package session

import (
	"context"
	"gonetlib/logger"
	"gonetlib/util"
	"net"
	"sync"
	"sync/atomic"
)

type ISessionManager interface {
	NewSession(uint64, net.Conn, ISessionHandler) ISession
	Dispose()
}

type SessionManager struct {
	ctx      context.Context
	pool     sync.Pool
	sessions sync.Map
	disposed int32
}

func CreateSessionManager(logger logger.ILogger, limit uint32, ctx context.Context) *SessionManager {
	return &SessionManager{
		ctx: ctx,
		pool: sync.Pool{
			New: func() interface{} {
				return newTcpSession(logger, ctx)
			},
		},
		sessions: sync.Map{},
		disposed: 0,
	}
}

func (s *SessionManager) NewSession(id uint64, conn net.Conn, handler ISessionHandler) ISession {
	if atomic.LoadInt32(&s.disposed) == 1 {
		return nil
	}

	session := s.pool.Get().(ISession)
	if session == nil {
		return nil
	}

	session.Setup(id, conn, handler, s)

	_, loaded := s.sessions.LoadOrStore(session.GetID(), session)
	if loaded == true {
		//말이 안되는 상황
		//풀에 있는 세션은 세션 관리 목록에 있으면 안됨
		return nil
	}

	return session
}

func (s *SessionManager) Dispose() {
	if util.InterlockedCompareExchange(&s.disposed, 1, 0) == false {
		return
	}

	s.sessions.Range(func(key, value any) bool {
		session := value.(ISession)
		session.Stop()
		return true
	})
}

func (s *SessionManager) OnRelease(session ISession) {
	if session == nil {
		return
	}

	//세션 관리 목록에서 삭제 후
	s.sessions.Delete(session.GetID())

	//세션 풀에 삽입
	s.pool.Put(session)
}
