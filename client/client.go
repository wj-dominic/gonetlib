package client

import (
	"gonetlib/logger"
	"gonetlib/session"
)

type IClientHandler interface {
	OnRun(logger.ILogger) error
	OnStop() error
	session.ISessionHandler
}

type IClient interface {
	Run() bool
	Stop() bool
}

type client struct {
	logger    logger.ILogger
	info      ClientInfo
	handler   IClientHandler
	connector IConnector
}

func newClient(logger logger.ILogger, info ClientInfo, handler IClientHandler) IClient {
	return &client{
		logger:    logger,
		info:      info,
		handler:   handler,
		connector: CreateConnector(logger, info.ServerAddress, info.Protocols),
	}
}

func (c *client) Run() bool {
	if err := c.connector.Start(); err != nil {
		c.logger.Error("Failed to start by connector", logger.Why("to", c.info.ServerAddress.ToString()), logger.Why("error", err.Error()))
		return false
	}

	if c.handler != nil {
		if err := c.handler.OnRun(c.logger); err != nil {
			c.logger.Error("Failed to call on run handler", logger.Why("error", err.Error()))
			return false
		}
	}

	c.logger.Info("Success to start client")
	return true
}

func (c *client) Stop() bool {
	if err := c.handler.OnStop(); err != nil {
		c.logger.Error("Failed to call on stop handler", logger.Why("error", err.Error()))
	}

	c.connector.Stop()
	c.logger.Dispose()

	return true
}
