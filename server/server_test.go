package server_test

import (
	"gonetlib/logger"
	"gonetlib/message"
	"gonetlib/server"
	"gonetlib/session"
	"gonetlib/util"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type EchoServer struct {
	logger logger.ILogger
	count  int32
}

func (s *EchoServer) OnRun(logger logger.ILogger) error {
	s.logger = logger
	s.count = 0
	return nil
}

func (s *EchoServer) OnStop() error {
	s.logger.Info("Connected session count", logger.Why("count", s.count))
	return nil
}

type EchoSession struct {
	EchoServer
}

func (h *EchoSession) Init(logger logger.ILogger) error {
	return nil
}

func (h *EchoSession) OnConnect(session session.ISession) error {
	h.logger.Info("On connect session", logger.Why("id", session.GetID()))
	util.InterlockIncrement(&h.count)
	return nil
}

func (h *EchoSession) OnRecv(session session.ISession, packet *message.Message) error {
	var msg string
	var id int
	packet.Pop(&msg)
	packet.Pop(&id)

	var sb strings.Builder
	sb.WriteString(msg)
	sb.WriteString(strconv.Itoa(id))

	h.logger.Info("On recv session", logger.Why("id", session.GetID()), logger.Why("msg", sb.String()))
	return nil
}

func (h *EchoSession) OnSend(session session.ISession, sentBytes []byte) error {
	return nil
}

func (h *EchoSession) OnDisconnect(session session.ISession) error {
	h.logger.Info("On disconnect session", logger.Why("id", session.GetID()))
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

	//time.Sleep(time.Second * 30)

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(number int) {
			defer func() {
				wg.Done()
			}()

			conn, err := net.Dial("tcp", "127.0.0.1:50000")
			if err != nil {
				return
			}

			defer func() {
				conn.Close()
			}()

			packet := message.NewMessage()
			packet.Push("hello my name is ")
			packet.Push(number)
			packet.MakeHeader()

			conn.Write(packet.GetBuffer())

			time.Sleep(time.Second * 2)
		}(i)
	}

	wg.Wait()

	server.Stop()
}
