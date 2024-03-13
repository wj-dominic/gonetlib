package session

import "net"

type SessionConfig struct {
	id      uint64
	conn    net.Conn
	handler ISessionHandler
}
