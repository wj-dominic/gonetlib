package network

type Protocol byte

const (
	TCP Protocol = 1
	UDP          = iota << 1
)
