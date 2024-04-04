package session

import (
	"fmt"
	"gonetlib/logger"
	"gonetlib/util"
	"net"
	"sync"
	"sync/atomic"
)

type ISessionManager interface {
	NewSession(uint64, net.Conn, ISessionHandler) (ISession, error)
	Dispose()
}

type SessionManager struct {
	logger     logger.ILogger
	pool       sync.Pool
	sessions   sync.Map
	isDisposed int32
}

func CreateSessionManager(logger logger.ILogger, limit uint32) *SessionManager {
	return &SessionManager{
		logger: logger,
		pool: sync.Pool{
			New: func() interface{} {
				return NewTcpSession(logger)
			},
		},
		sessions:   sync.Map{},
		isDisposed: 0,
	}
}

func (s *SessionManager) NewSession(id uint64, conn net.Conn, handler ISessionHandler) (ISession, error) {
	if atomic.LoadInt32(&s.isDisposed) == 1 {
		return nil, fmt.Errorf("session manager was disposed")
	}

	session := s.pool.Get().(ISession)
	if session == nil {
		return nil, fmt.Errorf("failed to get session from pool")
	}

	session.Setup(id, conn, handler, s)

	value, loaded := s.sessions.LoadOrStore(session.GetID(), session)
	if loaded == true {
		//말이 안되는 상황
		//풀에 있는 세션은 세션 관리 목록에 있으면 안됨
		loadedSession := value.(ISession)
		return nil, fmt.Errorf("already session is running, session id:%d, loaded session:%d ", session.GetID(), loadedSession.GetID())
	}

	s.logger.Info("new session", logger.Why("id", session.GetID()))

	return session, nil
}

func (s *SessionManager) Dispose() {
	if util.InterlockedCompareExchange(&s.isDisposed, 1, 0) == false {
		return
	}

	s.sessions.Range(func(key, value any) bool {
		session := value.(ISession)
		session.Stop()
		return true
	})
}

func (s *SessionManager) OnRelease(sessionID uint64, session ISession) {
	if session == nil {
		return
	}

	s.logger.Info("release session", logger.Why("id", sessionID))

	//세션 관리 목록에서 삭제 후
	s.sessions.Delete(sessionID)

	//세션 풀에 삽입
	s.pool.Put(session)
}
