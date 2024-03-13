package server

import (
	"context"
	"gonetlib/logger"
	"gonetlib/session"

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
}

type Server struct {
	config   ServerConfig
	acceptor IAcceptor
	sessions session.ISessionManager
	handler  IServerHandler
	logger   logger.ILogger

	ctx    context.Context
	cancel context.CancelFunc
}

func newServerWithContext(config ServerConfig, ctx context.Context) IServer {
	_ctx, cancel := context.WithCancel(ctx)
	server := &Server{
		config:   config,
		acceptor: CreateAcceptor(ctx, config.Protocols, config.Address),
		sessions: session.CreateSessionManager(ctx, config.MaxSession),
		ctx:      _ctx,
		cancel:   cancel,
	}

	server.acceptor.SetHandler(server)

	return server
}

func newServer(config ServerConfig) IServer {
	return newServerWithContext(config, context.Background())
}

func (s *Server) Run() bool {
	if s.acceptor.StartAccept() == false {
		s.logger.Error("Failed to start accept", logger.Why("address", s.config.Address.ToString()))
		return false
	}

	return true
}

func (s *Server) Stop() bool {
	s.cancel()
	s.logger.Dispose()
	return true
}

func (s *Server) OnAccept(conn net.Conn) {
	session := s.sessions.NewSession(conn)
	if session == nil {
		s.logger.Error("Failed to create new session",
			logger.Why("local address", conn.LocalAddr().String()),
			logger.Why("remote", conn.RemoteAddr().String()))
		conn.Close()
		return
	}

	session.Start()
	//TODO : start 이후 핸들러 등록
}

func (s *Server) OnRecvFrom(client net.Addr, recvData []byte, recvBytes uint32) {
	//TODO : UDP 처리, 커넥트를 가지고 연결을 찾거나 만들도록
}
