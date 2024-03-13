package session

import "gonetlib/message"

type ISession interface {
	Start()
	Stop()
	GetID() uint64
}

type ISessionHandler interface {
	OnConnect()
	OnDisconnect()
	OnRecv(packet *message.Message)
	OnSend([]byte)
}

type Session struct {
	config SessionConfig
}

func (session *Session) Start() {

}

func (session *Session) Stop() {

}

func (session *Session) GetID() uint64 {
	return session.config.id
}
