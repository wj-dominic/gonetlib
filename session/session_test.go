package session_test

import (
	"gonetlib/logger"
	"gonetlib/message"
	"gonetlib/session"
	"gonetlib/util/snowflake"
	"net"
	"testing"
	"time"
)

type ServerSession struct {
	logger logger.ILogger
}

func (server *ServerSession) Init(logger logger.ILogger) error {
	server.logger = logger
	return nil
}

func (server *ServerSession) OnConnect(session session.ISession) error {
	server.logger.Info("on connect server", logger.Why("id", session.GetID()))
	return nil
}
func (server *ServerSession) OnDisconnect(session session.ISession) error {
	server.logger.Info("on disconnect server", logger.Why("id", session.GetID()))
	return nil
}
func (server *ServerSession) OnRecv(session session.ISession, packet *message.Message) error {
	var msg string
	packet.Pop(&msg)

	server.logger.Info("recv message from client", logger.Why("msg", msg))
	session.Send(msg)

	return nil
}
func (server *ServerSession) OnSend(session session.ISession, sentBytes []byte) error {
	server.logger.Info("sent message to client", logger.Why("sentBytes", len(sentBytes)))
	return nil
}

type ClientSession struct {
	logger logger.ILogger
}

func (client *ClientSession) Init(logger logger.ILogger) error {
	client.logger = logger
	return nil
}

func (client *ClientSession) OnConnect(session session.ISession) error {
	client.logger.Info("on connect client", logger.Why("id", session.GetID()))

	msg := "test"
	client.logger.Info("send a message to server", logger.Why("msg", msg))
	session.Send(msg)

	return nil
}
func (client *ClientSession) OnDisconnect(session session.ISession) error {
	client.logger.Info("on disconnect client", logger.Why("id", session.GetID()))
	return nil
}
func (client *ClientSession) OnRecv(session session.ISession, packet *message.Message) error {
	var msg string
	packet.Pop(&msg)
	client.logger.Info("recv message from server", logger.Why("msg", msg))

	session.Send(msg)
	return nil
}
func (client *ClientSession) OnSend(session session.ISession, sentBytes []byte) error {
	client.logger.Info("sent message to server", logger.Why("sentBytes", len(sentBytes)))
	return nil
}

func TestSession(t *testing.T) {
	config := logger.CreateLoggerConfig().
		WriteToConsole().
		WriteToFile(
			logger.WriteToFile{
				Filepath:        "./test_session.log",
				RollingInterval: logger.RollingIntervalDay,
			}).
		MinimumLevel(logger.DebugLevel).
		TickDuration(1000)
	_logger := config.CreateLogger()

	server, client := net.Pipe()

	sessionManager := session.CreateSessionManager(_logger, 1000)
	serverSession, _ := sessionManager.NewSession(snowflake.GenerateID(1), server, &ServerSession{})
	clientSession, _ := sessionManager.NewSession(snowflake.GenerateID(1), client, &ClientSession{})

	if err := serverSession.Start(); err != nil {
		t.Error(err)
		t.Failed()
	}

	if err := clientSession.Start(); err != nil {
		t.Error(err)
		t.Failed()
	}

	time.Sleep(time.Second * 20)

	sessionManager.Dispose()
	_logger.Dispose()
}
