package client

import (
	"gonetlib/logger"
	"gonetlib/util/network"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var ReconnectDuration time.Duration = time.Second * 3
var ReconnectLimit uint32 = 0 //zero is unlimited

type ConnectHandler interface {
	OnConnect(network.Protocol, net.Conn)
}

type Connector interface {
	Start() error
	Stop()
}

type connector struct {
	logger         logger.ILogger
	serverAddress  network.Endpoint
	wg             sync.WaitGroup
	reconnect      chan bool
	reconnectCount uint32
	handler        ConnectHandler
	quit           atomic.Bool
}

type tcpConnector struct {
	connector
}

func NewTcpConnector(logger logger.ILogger, serverAddress network.Endpoint, handler ConnectHandler) Connector {
	return &tcpConnector{
		connector: connector{
			logger:         logger,
			serverAddress:  serverAddress,
			wg:             sync.WaitGroup{},
			reconnect:      make(chan bool),
			reconnectCount: 0,
			handler:        handler,
			quit:           atomic.Bool{},
		},
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

		if ReconnectLimit != 0 && tcpConn.reconnectCount >= ReconnectLimit {
			tcpConn.logger.Info("Cannot reconnect", logger.Why("reconnectCount", tcpConn.reconnectCount), logger.Why("reconnectLimit", ReconnectLimit))
			break
		}

		conn, err := net.Dial("tcp", tcpConn.serverAddress.ToString())
		if err != nil {
			tcpConn.logger.Error("Failed to connect tcp", logger.Why("serverAddress", tcpConn.serverAddress.ToString()), logger.Why("error", err.Error()))
			tcpConn.logger.Info("Try to reconnect..", logger.Why("serverAddress", tcpConn.serverAddress.ToString()))
			time.Sleep(ReconnectDuration)
			if ReconnectLimit != 0 {
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
