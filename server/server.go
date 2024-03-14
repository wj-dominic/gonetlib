package server

import (
	"context"
	"gonetlib/logger"
	"gonetlib/message"
	"gonetlib/session"
	"time"

	"net"
)

type IServerHandler interface {
	OnConnect()
	OnRecv(packet *message.Message)
	OnSend(sendBytes []byte)
	OnDisconnect()
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

func newServerWithContext(info ServerInfo, ctx context.Context) IServer {
	_ctx, cancel := context.WithCancel(ctx)
	server := &Server{
		info:     info,
		acceptor: CreateAcceptor(ctx, info.Protocols, info.Address),
		sessions: session.CreateSessionManager(ctx, info.MaxSession),
		ctx:      _ctx,
		cancel:   cancel,
	}

	server.acceptor.SetHandler(server)

	return server
}

func newServer(config ServerInfo) IServer {
	return newServerWithContext(config, context.Background())
}

func (s *Server) Run() bool {
	if s.acceptor.StartAccept() == false {
		s.logger.Error("Failed to start accept", logger.Why("address", s.info.Address.ToString()))
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
	session := s.sessions.NewSession(s.makeSessionId(), conn, s.handler)
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

func (s *Server) makeSessionId() uint64 {
	var sessionId uint64

	now := time.Now()
	sessionId = sessionId | uint64(s.info.Id)<<32
	sessionId = uint64(uint32(sessionId) | uint32(now.Unix()))

	return sessionId
}

func (s *Server) OnRecvFrom(client net.Addr, recvData []byte, recvBytes uint32) {
	//TODO : UDP 처리, 커넥트를 가지고 연결을 찾거나 만들도록
}
