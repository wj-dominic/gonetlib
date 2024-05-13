package client

import (
	"github.com/wj-dominic/gonetlib/util/network"
)

type ClientInfo struct {
	ServerAddress network.Endpoint
	Protocols     network.Protocol
	ConnectorInfo
}
