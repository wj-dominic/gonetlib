package mmo_server

import "gonetlib/session"

//TODO:node 세션 1:1 매핑 작업 필요
type Node struct {
	session session.ISession
}

func CreateNode(session session.ISession) *Node {
	return &Node{
		session: session,
	}
}
