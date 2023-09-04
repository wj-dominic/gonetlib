package gonet

import (
	"context"
	"net"
)

type IServerHandler interface {
	OnConnect()
	OnRecv([]byte)
	OnSend(uint32)
	OnDisconnect()
}

type IServer interface {
	Run() bool
	Stop() bool
	RegistHandler(IServerHandler)
}

type ISession interface {
	Start()
}

type ISessionManager interface {
	NewSession(context.Context, net.Conn) ISession
}

type Server struct {
	meta     *ServerConfig
	acceptor IAcceptor
	sessions ISessionManager
	handler  IServerHandler

	ctx    context.Context
	cancel context.CancelFunc
}

func NewServer(config func(*ServerConfig)) IServer {
	meta := &ServerConfig{
		Address:    Endpoint{IP: "127.0.0.1", Port: 50000},
		MaxSession: 100,
		Protocols:  TCP,
	}

	config(meta)

	ctx, cancel := context.WithCancel(context.Background())

	server := &Server{
		meta:     meta,
		acceptor: NewAcceptor(ctx, meta.Protocols, meta.Address),
		sessions: NewSessionManager(ctx, meta.MaxSession),
		ctx:      ctx,
		cancel:   cancel,
	}

	acceptEvent := AcceptEvent{
		OnAccept:   server.OnAccept,
		OnRecvFrom: server.OnRecvFrom,
	}

	server.acceptor.SetEvent(&acceptEvent)

	return server
}

func (s *Server) RegistHandler(handler IServerHandler) {
	s.handler = handler
}

func (s *Server) Run() bool {
	s.acceptor.StartAccept()
	return true
}

func (s *Server) Stop() bool {
	s.cancel()
	return true
}

func (s *Server) OnAccept(conn net.Conn) {
	session := s.sessions.NewSession(s.ctx, conn)
	if session == nil {
		return
	}

	session.Start()
	//TODO : start 이후 핸들러 등록
}

func (s *Server) OnRecvFrom(client net.Addr, recvData []byte, recvBytes uint32) {
	//TODO : UDP 처리, 커넥트를 가지고 연결을 찾거나 만들도록
}
