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
	builder := server.CreateServerBuilder()
	builder.Configuration(server.ServerConfig{
		Address:    server.Endpoint{IP: "0.0.0.0", Port: 50000},
		Protocols:  server.TCP | server.UDP,
		MaxSession: 10000,
	})
	builder.Logger(&logger.Logger{})
	builder.Handler(&EchoServerHandler{})

	server := builder.Build()
	server.Run()
}
