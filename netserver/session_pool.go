package netserver

import (
	"container/list"
	"errors"
	"fmt"
	. "gonetlib/netlogger"
	"gonetlib/session"
)

type sessionPool struct {
	pool    map[uint64]*session.Session
	idStack *list.List

	objCnt uint64
}

func NewSessionPool(initCnt uint64) *sessionPool {
	pool := &sessionPool{}

	pool.pool = make(map[uint64]*session.Session, initCnt)
	pool.idStack = list.New()

	for i := initCnt - 1; i > 0; i-- {
		pool.pool[i] = session.NewSession()
		pool.idStack.PushBack(i)
	}

	pool.objCnt = initCnt

	return pool
}

func (sp *sessionPool) acquireSession() (uint64, *session.Session) {
	if sp.idStack.Len() == 0 {
		// 스택, 맵 추가
		sp.increasePool()
	}

	sessIdElem := sp.idStack.Back()
	sessionId := sessIdElem.Value.(uint64)
	session := sp.pool[sessionId]
	if session != nil {
		sp.idStack.Remove(sessIdElem)
	}

	return sessionId, session
}

func (sp *sessionPool) releaseSession(sessionId uint64) error {
	session := sp.pool[sessionId]
	if session == nil {
		err := fmt.Sprintf("failed to find session(%d)", sessionId)
		GetLogger().Error(err)
		return errors.New(err)
	}

	session.Reset()
	sp.idStack.PushBack(sessionId)

	return nil
}

func (sp *sessionPool) increasePool() {
	for i := sp.objCnt * 2; i > sp.objCnt; i-- {
		sp.idStack.PushBack(i)
		sp.pool[i] = session.NewSession()
	}

	sp.objCnt *= 2
}

func (sp *sessionPool) getSession(sessionId uint64) *session.Session {
	return sp.pool[sessionId]
}

func (sp *sessionPool) getObjCount() uint64 {
	return sp.objCnt
}
