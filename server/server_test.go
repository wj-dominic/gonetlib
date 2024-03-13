package server_test

import (
	"gonetlib/logger"
	"gonetlib/server"
	"testing"
)

type EchoServerHandler struct {
}

func (h *EchoServerHandler) OnConnect() {

}

func (h *EchoServerHandler) OnRecv(recvData []byte) {

}

func (h *EchoServerHandler) OnSend(sendBytes uint32) {

}

func (h *EchoServerHandler) OnDisconnect() {

}

func TestMain(m *testing.T) {

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
	builder.Configuration(server.ServerConfig{
		Address:    server.Endpoint{IP: "0.0.0.0", Port: 50000},
		Protocols:  server.TCP | server.UDP,
		MaxSession: 10000,
	})
	builder.Logger(_logger)
	builder.Handler(&EchoServerHandler{})

	server := builder.Build()
	server.Run()
}
