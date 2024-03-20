package server

import (
	"gonetlib/logger"
	"net"
	"sync"
)

type IAcceptHandler interface {
	OnAccept(net.Conn)
	OnRecvFrom(net.Addr, []byte, uint32)
}

type IAcceptor interface {
	Start() error
	Stop()
}

const (
	MAX_BUFFER uint32 = 66535
)

type Acceptor struct {
	logger     logger.ILogger
	listener   net.Listener
	packetConn net.PacketConn

	protocols    Protocol
	endpoint     Endpoint
	listenConfig net.ListenConfig
	handler      IAcceptHandler

	wg sync.WaitGroup
}

func CreateAcceptor(logger logger.ILogger, protocols Protocol, endpoint Endpoint, handler IAcceptHandler) IAcceptor {
	return &Acceptor{
		logger:     logger,
		listener:   nil,
		packetConn: nil,

		protocols:    protocols,
		endpoint:     endpoint,
		listenConfig: net.ListenConfig{},
		handler:      handler,

		wg: sync.WaitGroup{},
	}
}

func (a *Acceptor) Start() error {
	if (a.protocols & TCP) == TCP {
		listener, err := net.Listen("tcp", a.endpoint.ToString())
		if err != nil {
			return err
		}

		a.listener = listener

		a.wg.Add(1)
		go a.waitForTCPConn()
	}

	if (a.protocols & UDP) == UDP {
		conn, err := net.ListenPacket("udp", a.endpoint.ToString())
		if err != nil {
			return err
		}

		a.packetConn = conn

		a.wg.Add(1)
		go a.waitForUDPConn()
	}

	return nil
}

func (a *Acceptor) waitForTCPConn() {
	defer a.wg.Done()

	for {
		conn, err := a.listener.Accept()
		if err != nil {
			a.logger.Error("Failed to accept for tcp connection", logger.Why("error", err.Error()))
			return
		}

		a.onAccept(conn)
	}
}

func (a *Acceptor) waitForUDPConn() {
	defer a.wg.Done()

	for {
		buffer := make([]byte, MAX_BUFFER)
		recvBytes, addr, err := a.packetConn.ReadFrom(buffer)
		if err != nil {
			a.logger.Error("Failed to read from buffer", logger.Why("error", err.Error()))
			return
		}

		a.onRecvFrom(addr, buffer, uint32(recvBytes))
	}
}

func (a *Acceptor) onAccept(conn net.Conn) {
	if a.handler != nil {
		a.handler.OnAccept(conn)
	}
}

func (a *Acceptor) onRecvFrom(client net.Addr, recvData []byte, recvBytes uint32) {
	if a.handler != nil {
		a.handler.OnRecvFrom(client, recvData, recvBytes)
	}
}

func (a *Acceptor) Stop() {
	if a.listener != nil {
		a.listener.Close()
	}

	if a.packetConn != nil {
		a.packetConn.Close()
	}

	a.wg.Wait()
}
