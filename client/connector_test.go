package client_test

import (
	"net"
	"testing"
	"time"

	"github.com/wj-dominic/gonetlib/client"
	"github.com/wj-dominic/gonetlib/logger"
	"github.com/wj-dominic/gonetlib/util/network"
)

type ConnectHandler struct {
}

func (handler *ConnectHandler) OnConnect(protocol network.Protocol, conn net.Conn) {

}

func TestConnector(t *testing.T) {
	connector := client.NewTcpConnector(logger.Default(), network.Endpoint{IP: "127.0.0.1", Port: 50000}, &ConnectHandler{}, client.DefaultConnectorInfo())
	connector.Start()

	time.Sleep(time.Second * 5)

	connector.Stop()
}
