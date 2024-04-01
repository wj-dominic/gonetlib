package client_test

import (
	"gonetlib/client"
	"gonetlib/logger"
	"gonetlib/message"
	"gonetlib/session"
	"gonetlib/util/network"
	"testing"
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
	return nil
}
func (echo *EchoClient) OnDisconnect(session session.ISession) error {
	return nil
}
func (echo *EchoClient) OnRecv(session session.ISession, packet *message.Message) error {
	return nil
}
func (echo *EchoClient) OnSend(session session.ISession, sentBytes []byte) error {
	return nil
}

func TestClient(t *testing.T) {
	config := logger.CreateLoggerConfig().
		WriteToConsole().
		WriteToFile(
			logger.WriteToFile{
				Filepath: "./EchoClient.log",
			}).
		MinimumLevel(logger.DebugLevel).
		TickDuration(1000)
	_logger := config.CreateLogger()

	builder := client.CreateClientBuilder()
	builder.Configuration(client.ClientInfo{
		ServerAddress: network.Endpoint{IP: "127.0.0.1", Port: 50000},
		Protocols:     network.TCP | network.UDP,
	})
	builder.Logger(_logger)
	builder.Handler(&EchoClient{})

	client := builder.Build()
	client.Run()
}
