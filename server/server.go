package server

import (
	"context"
	"gonetlib/logger"
	"gonetlib/session"
	"time"

	"net"
)

type IServerHandler interface {
	OnRun() error
	OnStop() error
	session.ISessionHandler
}

type IServer interface {
	Run() bool
	Stop() bool
}

type Server struct {
	info     ServerInfo
	acceptor IAcceptor
	sessions session.ISessionManager
	handler  IServerHandler
	logger   logger.ILogger

	ctx    context.Context
	cancel context.CancelFunc
}

func newServerWithContext(logger logger.ILogger, info ServerInfo, handler IServerHandler, ctx context.Context) IServer {
	_ctx, cancel := context.WithCancel(ctx)
	server := &Server{
		info:     info,
		acceptor: nil,
		sessions: session.CreateSessionManager(logger, info.MaxSession, ctx),
		handler:  handler,
		logger:   logger,

		ctx:    _ctx,
		cancel: cancel,
	}

	acceptor := CreateAcceptor(logger, info.Protocols, info.Address, server, ctx)
	server.acceptor = acceptor

	return server
}

func newServer(logger logger.ILogger, info ServerInfo, handler IServerHandler) IServer {
	return newServerWithContext(logger, info, handler, context.Background())
}

func (s *Server) Run() bool {
	if err := s.acceptor.Start(); err != nil {
		s.logger.Error("Failed to start accept",
			logger.Why("address", s.info.Address.ToString()),
			logger.Why("error", err.Error()))
		return false
	}

	if s.handler != nil {
		if err := s.handler.OnRun(); err != nil {
			s.logger.Error("Failed to call on run handler", logger.Why("error", err.Error()))
			return false
		}
	}

	s.logger.Info("Success to run server")
	return true
}

func (s *Server) Stop() bool {
	if err := s.handler.OnStop(); err != nil {
		s.logger.Error("Failed to call on stop handler", logger.Why("error", err.Error()))
		return false
	}

	//기능 중단
	s.cancel()

	//종료 대기
	s.acceptor.Stop()
	s.logger.Dispose()

	s.logger.Info("Success to stop server")
	return true
}

func (s *Server) OnAccept(conn net.Conn) {
	session := s.sessions.NewSession(s.makeSessionId(), conn, s.handler)
	if session == nil {
		s.logger.Error("Failed to create new session",
			logger.Why("local address", conn.LocalAddr().String()),
			logger.Why("remote", conn.RemoteAddr().String()))

		conn.Close()
		return
	}

	if err := session.Start(); err != nil {
		s.logger.Error("Failed to start a session",
			logger.Why("session id", session.GetID()),
			logger.Why("internal", err.Error()))

		conn.Close()
		return
	}
}

func (s *Server) makeSessionId() uint64 {
	//TODO:snowflake util 만들어서 이거 호출하게 하기

	var sessionId uint64

	now := time.Now()
	sessionId = uint64(s.info.Id) << 32
	sessionId |= uint64(uint32(now.Unix()))

	return sessionId
}

func (s *Server) OnRecvFrom(client net.Addr, recvData []byte, recvBytes uint32) {
	//TODO : UDP 처리, 커넥트를 가지고 연결을 찾거나 만들도록
}
