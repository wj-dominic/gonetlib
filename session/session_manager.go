package session

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/wj-dominic/gonetlib/logger"
	"github.com/wj-dominic/gonetlib/util"
)

type SessionManager interface {
	NewSession(uint64, net.Conn, SessionHandler) (Session, error)
	Dispose()
}

type sessionManager struct {
	logger     logger.Logger
	pool       sync.Pool
	sessions   sync.Map
	isDisposed int32
}

func NewSessionManager(logger logger.Logger, limit uint32) *sessionManager {
	return &sessionManager{
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

func (s *sessionManager) NewSession(id uint64, conn net.Conn, handler SessionHandler) (Session, error) {
	if atomic.LoadInt32(&s.isDisposed) == 1 {
		return nil, fmt.Errorf("session manager was disposed")
	}

	session := s.pool.Get().(Session)
	if session == nil {
		return nil, fmt.Errorf("failed to get session from pool")
	}

	session.Setup(id, conn, handler, s)

	value, loaded := s.sessions.LoadOrStore(session.GetID(), session)
	if loaded == true {
		//말이 안되는 상황
		//풀에 있는 세션은 세션 관리 목록에 있으면 안됨
		loadedSession := value.(Session)
		return nil, fmt.Errorf("already session is running, session id:%d, loaded session:%d ", session.GetID(), loadedSession.GetID())
	}

	s.logger.Info("new session", logger.Why("id", session.GetID()))

	return session, nil
}

func (s *sessionManager) Dispose() {
	if util.InterlockedCompareExchange(&s.isDisposed, 1, 0) == false {
		return
	}

	s.sessions.Range(func(key, value any) bool {
		session := value.(Session)
		session.Stop()
		return true
	})
}

func (s *sessionManager) OnRelease(sessionID uint64, session Session) {
	if session == nil {
		return
	}

	s.logger.Info("release session", logger.Why("id", sessionID))

	//세션 관리 목록에서 삭제 후
	s.sessions.Delete(sessionID)

	//세션 풀에 삽입
	s.pool.Put(session)
}
