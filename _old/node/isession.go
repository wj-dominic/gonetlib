package node

import "gonetlib/message"

type ISession interface {
	SendPost(packet *message.Message) bool
}