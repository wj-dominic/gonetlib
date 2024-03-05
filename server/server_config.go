package server

type Protocol byte

const (
	TCP Protocol = iota + 1
	UDP
)

type ServerConfig struct {
	Address   Endpoint
	Protocols Protocol

	MaxSession uint32
}
