package server_test

import (
	"gonetlib/logger"
	"gonetlib/message"
	"gonetlib/server"
	"gonetlib/session"
	"gonetlib/util"
	"gonetlib/util/network"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type TestServer struct {
	logger logger.Logger
	count  int32
}

func (s *TestServer) OnRun(logger logger.Logger) error {
	s.logger = logger
	s.count = 0
	return nil
}

func (s *TestServer) OnStop() error {
	s.logger.Info("Connected session count", logger.Why("count", s.count))
	return nil
}

func (s *TestServer) OnConnect(session session.Session) error {
	s.logger.Info("On connect session", logger.Why("id", session.GetID()))
	util.InterlockIncrement(&s.count)
	return nil
}

func (s *TestServer) OnRecv(session session.Session, packet *message.Message) error {
	var msg string
	var id int
	packet.Pop(&msg)
	packet.Pop(&id)

	var sb strings.Builder
	sb.WriteString(msg)
	sb.WriteString(strconv.Itoa(id))

	s.logger.Info("On recv session", logger.Why("id", session.GetID()), logger.Why("msg", sb.String()))
	return nil
}

func (s *TestServer) OnSend(session session.Session, sentBytes []byte) error {
	return nil
}

func (s *TestServer) OnDisconnect(session session.Session) error {
	s.logger.Info("On disconnect session", logger.Why("id", session.GetID()))
	return nil
}

func TestSever(t *testing.T) {
	config := logger.NewLoggerConfig().
		WriteToConsole().
		WriteToFile(
			logger.WriteToFile{
				Filepath:        "./EchoServer.log",
				RollingInterval: logger.RollingIntervalDay,
			}).
		MinimumLevel(logger.DebugLevel).
		TickDuration(1000)
	_logger := config.CreateLogger()

	builder := server.NewServerBuilder()
	builder.Configuration(server.ServerInfo{
		Id:         1,
		Address:    network.Endpoint{IP: "0.0.0.0", Port: 50000},
		Protocols:  network.TCP | network.UDP,
		MaxSession: 10000,
	})
	builder.Logger(_logger)
	builder.Handler(&TestServer{})

	server := builder.Build()
	server.Run()

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
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

		time.Sleep(time.Millisecond)
	}

	wg.Wait()

	server.Stop()
}
