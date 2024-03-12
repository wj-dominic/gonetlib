package server

import (
	"context"
	"gonetlib/logger"
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

type ISession interface {
	Start()
}

type ISessionManager interface {
	NewSession(context.Context, net.Conn) ISession
}

type Server struct {
	config   ServerConfig
	acceptor IAcceptor
	sessions ISessionManager
	handler  IServerHandler
	logger   logger.ILogger

	ctx    context.Context
	cancel context.CancelFunc
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
	session := s.sessions.NewSession(s.ctx, conn)
	if session == nil {
		s.logger.Error("Failed to create new session",
			logger.Why("local address", conn.LocalAddr().String()),
			logger.Why("remote", conn.RemoteAddr().String()))
		return
	}

	session.Start()
	//TODO : start 이후 핸들러 등록
}

func (s *Server) OnRecvFrom(client net.Addr, recvData []byte, recvBytes uint32) {
	//TODO : UDP 처리, 커넥트를 가지고 연결을 찾거나 만들도록
}
