package gonet

import (
	"context"
	"net"
	"sync"
)

type AcceptEvent struct {
	OnAccept   func(net.Conn)
	OnRecvFrom func(net.Addr, []byte, uint32)
}

type IAcceptor interface {
	StartAccept() bool
	StopAccept()
	SetEvent(*AcceptEvent)
}

const (
	MAX_BUFFER uint32 = 4096
)

type Acceptor struct {
	protocols    Protocol
	endpoint     Endpoint
	listenConfig net.ListenConfig
	ctx          context.Context
	wg           sync.WaitGroup
	event        *AcceptEvent
}

func NewAcceptor(ctx context.Context, protocols Protocol, endpoint Endpoint) IAcceptor {
	return &Acceptor{
		ctx:       ctx,
		protocols: protocols,
		endpoint:  endpoint,
	}
}

func (a *Acceptor) SetEvent(event *AcceptEvent) {
	a.event = event
}

func (a *Acceptor) StartAccept() bool {
	defer a.wg.Wait()

	if (a.protocols & TCP) == 1 {
		listener, err := a.listenConfig.Listen(a.ctx, "tcp", a.endpoint.ToString())
		if err != nil {
			return false
		}

		a.wg.Add(1)
		go func() {
			defer a.wg.Done()

			for {
				select {
				case <-a.ctx.Done():
					return
				default:
					conn, err := listener.Accept()
					if err != nil {
						return
					}

					a.wg.Add(1)
					go func(conn net.Conn) {
						defer a.wg.Done()

						a.onAccept(conn)

					}(conn)
				}
			}
		}()
	}

	if (a.protocols & UDP) == 1 {
		conn, err := a.listenConfig.ListenPacket(a.ctx, "udp", a.endpoint.ToString())
		if err != nil {
			return false
		}

		a.wg.Add(1)
		go func() {
			defer a.wg.Done()

			for {
				select {
				case <-a.ctx.Done():
					return
				default:
					buffer := make([]byte, MAX_BUFFER)
					recvBytes, addr, err := conn.ReadFrom(buffer)
					if err != nil {
						return
					}

					a.wg.Add(1)
					go func(buffer []byte, recvBytes int, addr net.Addr) {
						defer a.wg.Done()

						a.onRecvFrom(addr, buffer, uint32(recvBytes))

					}(buffer, recvBytes, addr)
				}
			}
		}()
	}

	return true
}

func (a *Acceptor) onAccept(conn net.Conn) {
	if a.event != nil {
		a.event.OnAccept(conn)
	}
}

func (a *Acceptor) onRecvFrom(client net.Addr, recvData []byte, recvBytes uint32) {
	if a.event != nil {
		a.event.OnRecvFrom(client, recvData, recvBytes)
	}
}

func (a *Acceptor) StopAccept() {
	_, cancel := context.WithCancel(a.ctx)
	cancel()
}
