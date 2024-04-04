package client_test

import (
	"gonetlib/client"
	"gonetlib/logger"
	"gonetlib/util/network"
	"net"
	"testing"
	"time"
)

type ConnectHandler struct {
}

func (handler *ConnectHandler) OnConnect(protocol network.Protocol, conn net.Conn) {

}

func TestConnector(t *testing.T) {
	connector := client.NewTcpConnector(logger.Default(), network.Endpoint{IP: "127.0.0.1", Port: 50000}, &ConnectHandler{})
	connector.Start()

	time.Sleep(time.Second * 5)

	connector.Stop()
}
