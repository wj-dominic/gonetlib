package client

import (
	"gonetlib/logger"
	"gonetlib/session"
	"gonetlib/util/network"
	"gonetlib/util/snowflake"
	"net"
	"sync"
	"sync/atomic"
)

type Client interface {
	Run() bool
	Stop() bool
}

type gonetClient struct {
	logger     logger.Logger
	info       ClientInfo
	handler    ClientHandler
	connectors map[network.Protocol]Connector
	sessions   sync.Map
	isStopped  atomic.Bool
}

func newClient(logger logger.Logger, info ClientInfo, handler ClientHandler) Client {
	client := &gonetClient{
		logger:     logger,
		info:       info,
		handler:    handler,
		connectors: make(map[network.Protocol]Connector),
		sessions:   sync.Map{},
		isStopped:  atomic.Bool{},
	}

	if network.IsTCP(info.Protocols) == true {
		client.connectors[network.TCP] = NewTcpConnector(logger, info.ServerAddress, client, info.ConnectorInfo)
	}

	if network.IsUDP(info.Protocols) == true {
		client.connectors[network.UDP] = NewUdpConnector(logger, info.ServerAddress, client, info.ConnectorInfo)
	}

	return client
}

func (c *gonetClient) Run() bool {
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

	c.logger.Info("Success to run client")
	return true
}

func (c *gonetClient) Stop() bool {
	if c.isStopped.CompareAndSwap(false, true) == false {
		return false
	}

	if err := c.handler.OnStop(); err != nil {
		c.logger.Error("Failed to call on stop handler", logger.Why("error", err.Error()))
	}

	for _, connector := range c.connectors {
		connector.Stop()
	}

	c.sessions.Range(func(key, value any) bool {
		session := value.(struct {
			network.Protocol
			session.Session
		})
		session.Session.Stop()
		return true
	})

	c.logger.Dispose()

	return true
}

func (c *gonetClient) OnConnect(protocol network.Protocol, conn net.Conn) {
	if network.IsTCP(protocol) == true {
		tcpSession := session.NewTcpSession(c.logger)
		tcpSession.Setup(snowflake.GenerateID(1), conn, c.handler, c)

		c.sessions.Store(tcpSession.GetID(), struct {
			network.Protocol
			session.Session
		}{network.TCP, tcpSession})

		tcpSession.Start()
	}

	if network.IsUDP(protocol) == true {
		//TODO:UDP 세션 작업
	}
}

func (c *gonetClient) OnRelease(id uint64, inSession session.Session) {
	value, loaded := c.sessions.LoadAndDelete(id)
	if loaded == false {
		c.logger.Error("Failed to load from client sessions", logger.Why("id", id))
		return
	}

	foundSession := value.(struct {
		network.Protocol
		session.Session
	})

	if foundSession.GetID() != inSession.GetID() {
		c.logger.Error("Not match session id between found and in params", logger.Why("found", foundSession.GetID()), logger.Why("in", inSession.GetID()))
		panic(1)
	}

	if c.isStopped.Load() == false {
		c.connectors[foundSession.Protocol].Reconnect()
	}
}
