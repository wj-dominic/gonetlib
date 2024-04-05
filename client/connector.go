package client

import (
	"gonetlib/logger"
	"gonetlib/util/network"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type ConnectHandler interface {
	OnConnect(network.Protocol, net.Conn)
}

type Connector interface {
	Start() error
	Stop()
	Reconnect()
}

type ConnectorInfo struct {
	reconnectDuration time.Duration
	reconnectLimit    uint32
}

func DefaultConnectorInfo() ConnectorInfo {
	return NewConnectorInfo(time.Second*3, 0)
}

func NewConnectorInfo(reconnectDuration time.Duration, reconnectLimit uint32) ConnectorInfo {
	return ConnectorInfo{
		reconnectDuration: reconnectDuration,
		reconnectLimit:    reconnectLimit,
	}
}

type connector struct {
	logger         logger.ILogger
	serverAddress  network.Endpoint
	wg             sync.WaitGroup
	reconnect      chan bool
	reconnectCount uint32
	option         ConnectorInfo
	handler        ConnectHandler
	quit           atomic.Bool
}

func newConnector(logger logger.ILogger, serverAddress network.Endpoint, handler ConnectHandler, option ConnectorInfo) connector {
	return connector{
		logger:         logger,
		serverAddress:  serverAddress,
		wg:             sync.WaitGroup{},
		reconnect:      make(chan bool),
		reconnectCount: 0,
		option:         option,
		handler:        handler,
		quit:           atomic.Bool{},
	}
}

func (c *connector) Reconnect() {
	c.reconnect <- true
}

type tcpConnector struct {
	connector
}

func NewTcpConnector(logger logger.ILogger, serverAddress network.Endpoint, handler ConnectHandler, option ConnectorInfo) Connector {
	return &tcpConnector{
		connector: newConnector(logger, serverAddress, handler, option),
	}
}

func (tcpConn *tcpConnector) Start() error {
	tcpConn.wg.Add(1)
	go tcpConn.tryConnect()
	return nil
}

func (tcpConn *tcpConnector) tryConnect() {
	defer tcpConn.wg.Done()

	for {
		if tcpConn.quit.Load() == true {
			break
		}

		if tcpConn.option.reconnectLimit != 0 && tcpConn.reconnectCount >= tcpConn.option.reconnectLimit {
			tcpConn.logger.Info("Cannot reconnect", logger.Why("reconnectCount", tcpConn.reconnectCount), logger.Why("reconnectLimit", tcpConn.option.reconnectLimit))
			break
		}

		conn, err := net.Dial("tcp", tcpConn.serverAddress.ToString())
		if err != nil {
			tcpConn.logger.Error("Failed to connect tcp", logger.Why("serverAddress", tcpConn.serverAddress.ToString()), logger.Why("error", err.Error()))
			tcpConn.logger.Info("Try to reconnect..", logger.Why("serverAddress", tcpConn.serverAddress.ToString()))
			time.Sleep(tcpConn.option.reconnectDuration)
			if tcpConn.option.reconnectLimit != 0 {
				tcpConn.reconnectCount++
			}
			continue
		}

		tcpConn.onConnect(conn)

		shouldReconnect, ok := <-tcpConn.reconnect
		if ok == false || shouldReconnect == false {
			break
		}
	}
}

func (tcpConn *tcpConnector) onConnect(conn net.Conn) {
	if tcpConn.handler != nil {
		tcpConn.handler.OnConnect(network.TCP, conn)
	}
}

func (tcpConn *tcpConnector) Stop() {
	tcpConn.quit.Store(true)
	close(tcpConn.reconnect)
	tcpConn.wg.Wait()
}

type udpConnector struct {
	connector
}

func NewUdpConnector(logger logger.ILogger, serverAddress network.Endpoint, handler ConnectHandler, option ConnectorInfo) Connector {
	return &udpConnector{
		connector: newConnector(logger, serverAddress, handler, option),
	}
}

func (udpConn *udpConnector) Start() error {
	udpConn.wg.Add(1)
	go udpConn.tryConnect()
	return nil
}

func (udpConn *udpConnector) tryConnect() {
	defer udpConn.wg.Done()

	for {
		if udpConn.quit.Load() == true {
			break
		}

		if udpConn.option.reconnectLimit != 0 && udpConn.reconnectCount >= udpConn.option.reconnectLimit {
			udpConn.logger.Info("Cannot reconnect", logger.Why("reconnectCount", udpConn.reconnectCount), logger.Why("reconnectLimit", udpConn.option.reconnectLimit))
			break
		}

		conn, err := net.Dial("udp", udpConn.serverAddress.ToString())
		if err != nil {
			udpConn.logger.Error("Failed to connect udp", logger.Why("serverAddress", udpConn.serverAddress.ToString()), logger.Why("error", err.Error()))
			udpConn.logger.Info("Try to reconnect..", logger.Why("serverAddress", udpConn.serverAddress.ToString()))
			time.Sleep(udpConn.option.reconnectDuration)
			if udpConn.option.reconnectLimit != 0 {
				udpConn.reconnectCount++
			}
			continue
		}

		udpConn.onConnect(conn)

		shouldReconnect, ok := <-udpConn.reconnect
		if ok == false || shouldReconnect == false {
			break
		}
	}
}

func (udpConn *udpConnector) onConnect(conn net.Conn) {
	if udpConn.handler != nil {
		udpConn.handler.OnConnect(network.UDP, conn)
	}
}

func (udpConn *udpConnector) Stop() {
	udpConn.quit.Store(true)
	close(udpConn.reconnect)
	udpConn.wg.Wait()
}
