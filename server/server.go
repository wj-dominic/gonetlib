package server

import (
	"gonetlib/logger"
	"gonetlib/session"
	"gonetlib/util/snowflake"

	"net"
)

type Server interface {
	Run() bool
	Stop() bool
}

type gonetServer struct {
	info     ServerInfo
	acceptor Acceptor
	sessions session.SessionManager
	handler  ServerHandler
	logger   logger.Logger
}

func newServer(logger logger.Logger, info ServerInfo, handler ServerHandler) Server {
	server := &gonetServer{
		info:     info,
		acceptor: nil,
		sessions: session.NewSessionManager(logger, info.MaxSession),
		handler:  handler,
		logger:   logger,
	}

	acceptor := NewAcceptor(logger, info.Protocols, info.Address, server)
	server.acceptor = acceptor

	return server
}

func (s *gonetServer) Run() bool {
	if err := s.acceptor.Start(); err != nil {
		s.logger.Error("Failed to start accept",
			logger.Why("address", s.info.Address.ToString()),
			logger.Why("error", err.Error()))
		return false
	}

	if s.handler != nil {
		if err := s.handler.OnRun(s.logger); err != nil {
			s.logger.Error("Failed to call on run handler", logger.Why("error", err.Error()))
			return false
		}
	}

	s.logger.Info("Success to run server")
	return true
}

func (s *gonetServer) Stop() bool {
	if err := s.handler.OnStop(); err != nil {
		s.logger.Error("Failed to call on stop handler", logger.Why("error", err.Error()))
		return false
	}

	//종료 대기
	s.acceptor.Stop()
	s.sessions.Dispose()
	s.logger.Dispose()

	s.logger.Info("Success to stop server")
	return true
}

func (s *gonetServer) OnAccept(conn net.Conn) {
	session, err := s.sessions.NewSession(s.makeSessionId(), conn, s.handler)
	if session == nil {
		s.logger.Error("Failed to create new session",
			logger.Why("local address", conn.LocalAddr().String()),
			logger.Why("remote", conn.RemoteAddr().String()),
			logger.Why("error", err.Error()))

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

func (s *gonetServer) makeSessionId() uint64 {
	return snowflake.GenerateID(int64(s.info.Id))
}

func (s *gonetServer) OnRecvFrom(client net.Addr, recvData []byte, recvBytes uint32) {
	//TODO : UDP 처리, 커넥트를 가지고 연결을 찾거나 만들도록
}
