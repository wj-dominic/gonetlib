package mmo_server

import "gonetlib/session"

type INode interface {
	Send(interface{})
	SetContext(IPacketContext)
	Wait()
	Clear()
}

type Node struct {
	session session.ISession
	context IPacketContext
}

func NewNode(session session.ISession) INode {
	return &Node{
		session: session,
	}
}

func (n *Node) Send(data interface{}) {
	if n.session != nil {
		n.session.Send(data)
	}
}

func (n *Node) SetContext(context IPacketContext) {
	n.context = context
}

func (n *Node) Wait() {
	if n.context != nil {
		n.context.Wait()
	}
}

func (n *Node) Clear() {
	n.session = nil
	n.context = nil
}
