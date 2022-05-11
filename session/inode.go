package session

import "gonetlib/message"

type INode interface {
	OnConnect()
	OnDisconnect()
	OnRecv(packet *message.Message) bool
	OnSend(sendBytes int)
}
