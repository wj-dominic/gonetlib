package session

import (
	"gonetlib/message"
	"net"
)

type ISession interface {
	Start() bool
	Stop() bool
	Setup(uint64, net.Conn, ISessionHandler)
	GetID() uint64
}

type ISessionHandler interface {
	OnConnect()
	OnDisconnect()
	OnRecv(packet *message.Message)
	OnSend([]byte)
}

type Session struct {
	id       uint64
	conn     net.Conn
	handler  ISessionHandler
	refCount int32
}

func (session *Session) Setup(id uint64, conn net.Conn, handler ISessionHandler) {
	session.id = id
	session.conn = conn
	session.handler = handler
}

func (session *Session) GetID() uint64 {
	return session.id
}
