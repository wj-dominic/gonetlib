package netserver

import (
	"errors"
	. "gonetlib/netlogger"
	"gonetlib/node"
	"net"
	"time"
)

const (
	initPoolCount uint64 = 5
)

type SessionManager struct {
	pool *sessionPool
	stop chan bool
}

func NewSessionManager() *SessionManager {
	sessionMgr := &SessionManager{
		pool: NewSessionPool(initPoolCount),
		stop: make(chan bool),
	}

	return sessionMgr
}

func (sm *SessionManager) Run() {
	go sm.checkSessionProc()
}

func (sm *SessionManager) Stop() {
	sm.stop <- true
}

func (sm *SessionManager) RequestNewSession(conn net.Conn /* conn, node */) error {
	sessionId, session := sm.pool.acquireSession()
	if session == nil {
		err := errors.New("failed to acquireSession()")
		GetLogger().Error(err.Error())
		return err
	}

	node := node.NewUserNode()

	session.Setup(sessionId, conn, node)
	session.Start()

	return nil
}

func (sm *SessionManager) checkSessionProc() {
	ticker := time.NewTicker(time.Second * 1)

	for {
		select {
		case <-sm.stop:
			return
		case <-ticker.C:
			sm.checkSession()
		}
	}
}

func (sm *SessionManager) checkSession() {
	for i := sm.pool.getObjCount(); i > 0; i++ {
		session := sm.pool.getSession(i)
		if session == nil {
			GetLogger().Error("session(%d) not found", i)
			continue
		}

		if session.IsActive {
			continue
		}

		sm.pool.releaseSession(i)
	}
}
