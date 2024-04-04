package client

import (
	"gonetlib/logger"
	"gonetlib/session"
	"gonetlib/util/network"
	"gonetlib/util/snowflake"
	"net"
)

type IClientHandler interface {
	OnRun(logger.ILogger) error
	OnStop() error
	session.ISessionHandler
}

type Client interface {
	Run() bool
	Stop() bool
}

// TODO:여러 서버들과 연결을 맺을 수 있도록
type client struct {
	logger     logger.ILogger
	info       ClientInfo
	handler    IClientHandler
	connectors map[network.Protocol]Connector
	sessions   map[network.Protocol]session.ISession
}

func newClient(logger logger.ILogger, info ClientInfo, handler IClientHandler) Client {
	client := &client{
		logger:     logger,
		info:       info,
		handler:    handler,
		connectors: make(map[network.Protocol]Connector),
		sessions:   make(map[network.Protocol]session.ISession),
	}

	if network.IsTCP(info.Protocols) == true {
		client.connectors[network.TCP] = NewTcpConnector(logger, info.ServerAddress, client)
	}

	//TODO:UDP 커넥터
	if network.IsUDP(info.Protocols) == true {
		client.connectors[network.UDP] = NewTcpConnector(logger, info.ServerAddress, client)
	}

	return client
}

func (c *client) Run() bool {
	for _, connector := range c.connectors {
		if err := connector.Start(); err != nil {
			c.logger.Error("Failed to start by connector", logger.Why("to", c.info.ServerAddress.ToString()), logger.Why("error", err.Error()))
			return false
		}
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

	for _, connector := range c.connectors {
		connector.Stop()
	}

	c.logger.Dispose()

	return true
}

func (c *client) OnConnect(protocol network.Protocol, conn net.Conn) {
	if network.IsTCP(protocol) == true {
		session := session.NewTcpSession(c.logger)
		session.Setup(snowflake.GenerateID(1), conn, c.handler, c)
		c.sessions[network.TCP] = session
	}

	if network.IsUDP(protocol) == true {

	}
}

func (c *client) OnRelease(id uint64, session session.ISession) {
	var protocol network.Protocol
	for key, session := range c.sessions {
		if session.GetID() == id {
			protocol = key
			break
		}
	}

}
