package gonet_test

import (
	"gonetlib/gonet"
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
	server := gonet.NewServer(func(config *gonet.ServerConfig) {
		config.Address = gonet.Endpoint{IP: "127.0.0.1", Port: 50000}
		config.Protocols = gonet.TCP | gonet.UDP
		config.MaxSession = 10000
	})

	server.RegistHandler(&EchoServerHandler{})

	server.Run()
}
