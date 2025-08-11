package server

import (
	"github.com/wj-dominic/gonetlib/util/network"
)

type ServerInfo struct {
	Id         uint16
	Address    network.Endpoint
	Protocols  network.Protocol
	MaxSession uint32
}
