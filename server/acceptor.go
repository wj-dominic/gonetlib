package server

import (
	"net"
	"sync"

	"github.com/wj-dominic/gonetlib/util/network"

	"github.com/wj-dominic/gonetlib/logger"
)

type AcceptHandler interface {
	OnAccept(net.Conn)
	OnRecvFrom(net.Addr, []byte, uint32)
}

type Acceptor interface {
	Start() error
	Stop()
}

const (
	MAX_BUFFER uint32 = 66535
)

type acceptor struct {
	logger     logger.Logger
	listener   net.Listener
	packetConn net.PacketConn

	protocols network.Protocol
	endpoint  network.Endpoint
	handler   AcceptHandler
	wg        sync.WaitGroup
}

func NewAcceptor(logger logger.Logger, protocols network.Protocol, endpoint network.Endpoint, handler AcceptHandler) Acceptor {
	return &acceptor{
		logger:     logger,
		listener:   nil,
		packetConn: nil,

		protocols: protocols,
		endpoint:  endpoint,
		handler:   handler,

		wg: sync.WaitGroup{},
	}
}

func (a *acceptor) Start() error {
	if (a.protocols & network.TCP) == network.TCP {
		listener, err := net.Listen("tcp", a.endpoint.ToString())
		if err != nil {
			return err
		}

		a.listener = listener

		a.wg.Add(1)
		go a.waitForTCPConn()
	}

	if (a.protocols & network.UDP) == network.UDP {
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

func (a *acceptor) waitForTCPConn() {
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

func (a *acceptor) waitForUDPConn() {
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

func (a *acceptor) onAccept(conn net.Conn) {
	if a.handler != nil {
		a.handler.OnAccept(conn)
	}
}

func (a *acceptor) onRecvFrom(client net.Addr, recvData []byte, recvBytes uint32) {
	if a.handler != nil {
		a.handler.OnRecvFrom(client, recvData, recvBytes)
	}
}

func (a *acceptor) Stop() {
	if a.listener != nil {
		a.listener.Close()
	}

	if a.packetConn != nil {
		a.packetConn.Close()
	}

	a.wg.Wait()
}
