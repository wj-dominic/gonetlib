package client

import (
	"gonetlib/logger"
	"gonetlib/util/network"
)

type IConnector interface {
	Start() error
	Stop()
}

type connector struct {
	logger        logger.ILogger
	serverAddress network.Endpoint
	protocols     network.Protocol
}

func CreateConnector(logger logger.ILogger, serverAddress network.Endpoint, protocols network.Protocol) IConnector {
	return &connector{
		logger:        logger,
		serverAddress: serverAddress,
		protocols:     protocols,
	}
}

func (c *connector) Start() error {
	return nil
}

func (c *connector) Stop() {

}
