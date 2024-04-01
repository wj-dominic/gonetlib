package client

import (
	"gonetlib/util/network"
)

type ClientInfo struct {
	ServerAddress network.Endpoint
	Protocols     network.Protocol
}
