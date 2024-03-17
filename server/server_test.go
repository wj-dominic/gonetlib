package server_test

import (
	"gonetlib/logger"
	"gonetlib/message"
	"gonetlib/server"
	"testing"
)

type EchoServer struct {
}

func (s *EchoServer) OnRun() error {
	return nil
}

func (s *EchoServer) OnStop() error {
	return nil
}

type EchoSession struct {
	EchoServer
}

func (h *EchoSession) OnConnect() error {
	return nil
}

func (h *EchoSession) OnRecv(packet *message.Message) error {
	return nil
}

func (h *EchoSession) OnSend(sendBytes []byte) error {
	return nil
}

func (h *EchoSession) OnDisconnect() error {
	return nil
}

func TestSever(t *testing.T) {
	config := logger.CreateLoggerConfig().
		WriteToConsole().
		WriteToFile(
			logger.WriteToFile{
				Filepath:        "./EchoServer.log",
				RollingInterval: logger.RollingIntervalDay,
			}).
		MinimumLevel(logger.DebugLevel).
		TickDuration(1000)
	_logger := config.CreateLogger()

	builder := server.CreateServerBuilder()
	builder.Configuration(server.ServerInfo{
		Id:         1,
		Address:    server.Endpoint{IP: "0.0.0.0", Port: 50000},
		Protocols:  server.TCP | server.UDP,
		MaxSession: 10000,
	})
	builder.Logger(_logger)
	builder.Handler(&EchoSession{})

	server := builder.Build()
	server.Run()
	server.Stop()
}
