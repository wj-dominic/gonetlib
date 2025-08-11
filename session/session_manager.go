package session

import (
	"fmt"
	"gonetlib/logger"
	"gonetlib/util"
	"net"
	"sync"
	"sync/atomic"
)

type SessionManager interface {
	NewSession(uint64, net.Conn, SessionHandler) (Session, error)
	Dispose()
}

type sessionManager struct {
	logger         logger.Logger
	pool           sync.Pool
	sessions       sync.Map
	monitoringData MonitoringData
	isDisposed     int32
}

func NewSessionManager(logger logger.Logger, limit uint32) *sessionManager {
	return &sessionManager{
		logger: logger,
		pool: sync.Pool{
			New: func() interface{} {
				return NewTcpSession(logger)
			},
		},
		sessions: sync.Map{},
		monitoringData: MonitoringData{
			BySession: make(map[uint64]SessionMonitoringData),
		},
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

func (s *sessionManager) Collect() (interface{}, error) {
	// 매번 초기화 방식으로 할지 고민(or NewSession, OnRelease에서 lock), session count를 관리하는 부분도 추가 고려
	s.monitoringData.BySession = make(map[uint64]SessionMonitoringData)
	s.monitoringData.ActiveSessions = 0

	s.sessions.Range(func(key, value any) bool {
		session := value.(Session)

		s.monitoringData.ActiveSessions++
		s.monitoringData.BySession[session.GetID()] = session.SessionMonitoringData()

		return true
	})

	return s.monitoringData, nil
}
