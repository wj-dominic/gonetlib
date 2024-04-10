package client_test

import (
	"gonetlib/client"
	"gonetlib/logger"
	"gonetlib/message"
	"gonetlib/server"
	"gonetlib/session"
	"gonetlib/util/network"
	"testing"
	"time"
)

type EchoClient struct {
	logger logger.ILogger
}

func (echo *EchoClient) OnRun(logger logger.ILogger) error {
	echo.logger = logger
	return nil
}

func (echo *EchoClient) OnStop() error {
	return nil
}

func (echo *EchoClient) OnConnect(session session.ISession) error {
	echo.logger.Info("On connected server", logger.Why("id", session.GetID()))
	msg := "hello my name is echo client"
	session.Send(msg)
	return nil
}
func (echo *EchoClient) OnDisconnect(session session.ISession) error {
	echo.logger.Info("On disconnected server", logger.Why("id", session.GetID()))
	return nil
}
func (echo *EchoClient) OnRecv(session session.ISession, packet *message.Message) error {
	var msg string
	packet.Pop(&msg)
	echo.logger.Info("On recv message from server", logger.Why("id", session.GetID()), logger.Why("msg", msg))

	newMsg := "thanks!"

	session.Send(newMsg)
	return nil
}
func (echo *EchoClient) OnSend(session session.ISession, sentBytes []byte) error {
	return nil
}

type EchoServer struct {
	logger logger.ILogger
}

func (echo *EchoServer) OnRun(logger logger.ILogger) error {
	echo.logger = logger
	return nil
}

func (echo *EchoServer) OnStop() error {
	return nil
}

func (echo *EchoServer) OnConnect(session session.ISession) error {
	echo.logger.Info("On connected client", logger.Why("id", session.GetID()))
	return nil
}
func (echo *EchoServer) OnDisconnect(session session.ISession) error {
	echo.logger.Info("On disconnected client", logger.Why("id", session.GetID()))
	return nil
}
func (echo *EchoServer) OnRecv(session session.ISession, packet *message.Message) error {
	var msg string
	packet.Pop(&msg)
	echo.logger.Info("On recv message from client", logger.Why("id", session.GetID()), logger.Why("msg", msg))
	session.Send(msg)
	return nil
}
func (echo *EchoServer) OnSend(session session.ISession, sentBytes []byte) error {
	return nil
}

func TestClient(t *testing.T) {
	config := logger.NewLoggerConfig().
		WriteToConsole().
		WriteToFile(
			logger.WriteToFile{
				Filepath: "./EchoClient.log",
			}).
		MinimumLevel(logger.DebugLevel)
	logger := config.CreateLogger()

	//server
	serverBuilder := server.NewServerBuilder()
	serverBuilder.Configuration(server.ServerInfo{
		Id:         1,
		Address:    network.Endpoint{IP: "0.0.0.0", Port: 50000},
		Protocols:  network.TCP | network.UDP,
		MaxSession: 1000,
	})

	serverBuilder.Logger(logger)
	serverBuilder.Handler(&EchoServer{})

	server := serverBuilder.Build()
	server.Run()

	//client
	builder := client.NewClientBuilder()
	builder.Configuration(client.ClientInfo{
		ServerAddress: network.Endpoint{IP: "127.0.0.1", Port: 50000},
		Protocols:     network.TCP | network.UDP,
	})
	builder.Logger(logger)
	builder.Handler(&EchoClient{})

	client := builder.Build()
	client.Run()

	time.Sleep(time.Second * 20)

	client.Stop()
	server.Stop()
}
