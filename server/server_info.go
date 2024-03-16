package server

type Protocol byte

const (
	TCP Protocol = 1
	UDP          = iota << 1
)

type ServerInfo struct {
	Id         uint16
	Address    Endpoint
	Protocols  Protocol
	MaxSession uint32
}
