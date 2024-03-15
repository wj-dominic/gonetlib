package server

import (
	"context"
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
	SetHandler(handler IAcceptHandler)
	Dispose()
}

const (
	MAX_BUFFER uint32 = 4096
)

type Acceptor struct {
	protocols    Protocol
	endpoint     Endpoint
	listenConfig net.ListenConfig
	handler      IAcceptHandler

	logger logger.ILogger
	ctx    context.Context
	wg     sync.WaitGroup
}

func CreateAcceptor(logger logger.ILogger, protocols Protocol, endpoint Endpoint, ctx context.Context) IAcceptor {
	return &Acceptor{
		ctx:       ctx,
		protocols: protocols,
		endpoint:  endpoint,
		logger:    logger,
	}
}

func (a *Acceptor) SetHandler(handler IAcceptHandler) {
	a.handler = handler
}

func (a *Acceptor) Start() error {
	if (a.protocols & TCP) == 1 {
		listener, err := a.listenConfig.Listen(a.ctx, "tcp", a.endpoint.ToString())
		if err != nil {
			return err
		}

		a.wg.Add(1)
		go a.waitForTCPConn(listener)
	}

	if (a.protocols & UDP) == 1 {
		conn, err := a.listenConfig.ListenPacket(a.ctx, "udp", a.endpoint.ToString())
		if err != nil {
			return err
		}

		a.wg.Add(1)
		go a.waitForUDPConn(conn)
	}

	return nil
}

func (a *Acceptor) waitForTCPConn(listener net.Listener) {
	defer a.wg.Done()

	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				a.logger.Error("Failed to accept for tcp connection", logger.Why("error", err.Error()))
				return
			}

			a.onAccept(conn)
		}
	}
}

func (a *Acceptor) waitForUDPConn(conn net.PacketConn) {
	defer a.wg.Done()

	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			buffer := make([]byte, MAX_BUFFER)
			recvBytes, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				a.logger.Error("Failed to read from buffer", logger.Why("error", err.Error()))
				return
			}

			a.onRecvFrom(addr, buffer, uint32(recvBytes))
		}
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
	_, cancel := context.WithCancel(a.ctx)
	cancel()
	a.Dispose()
}

func (a *Acceptor) Dispose() {
	a.wg.Wait()
}
