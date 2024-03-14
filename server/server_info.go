package server

type Protocol byte

const (
	TCP Protocol = iota + 1
	UDP
)

type ServerInfo struct {
	Id        uint16
	Address   Endpoint
	Protocols Protocol

	MaxSession uint32
}
