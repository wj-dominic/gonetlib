package server_test

import (
	"fmt"
	"gonetlib/logger"
	"gonetlib/message"
	"gonetlib/server"
	"gonetlib/session"
	"gonetlib/task"
	"gonetlib/util"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

type TestServer struct {
	logger logger.ILogger
	count  int32
}

func (s *TestServer) OnRun(logger logger.ILogger) error {
	s.logger = logger
	s.count = 0
	return nil
}

func (s *TestServer) OnStop() error {
	s.logger.Info("Connected session count", logger.Why("count", s.count))
	return nil
}

type TestSession struct {
	TestServer
}

func (h *TestSession) Init(logger logger.ILogger) error {
	return nil
}

func (h *TestSession) OnConnect(session session.ISession) error {
	h.logger.Info("On connect session", logger.Why("id", session.GetID()))
	util.InterlockIncrement(&h.count)
	return nil
}

func (h *TestSession) OnRecv(session session.ISession, packet *message.Message) error {
	var packetId uint16
	var payloadSize uint16

	headerSize := uint16(util.Sizeof(reflect.ValueOf(packetId)) + util.Sizeof(reflect.ValueOf(payloadSize)))

	packet.Peek(&packetId)
	packet.Peek(&payloadSize)

	if packet.GetPayloadSize() < headerSize+payloadSize {
		return fmt.Errorf("payload size grater then packet size, payloadSize %d > packetSize %d", headerSize+payloadSize, packet.GetPayloadSize())
	}

	packet.MoveFront(headerSize)

	if packet.GetPayloadSize() != payloadSize {
		return fmt.Errorf("not matched payload size with packet size, payloadSize %d > packetSize %d", payloadSize, packet.GetPayloadSize())
	}

	//TODO:packet id로 packet handler 뽑기
	handler := GetPacketHandler(packetId)
	if handler != nil {
		return fmt.Errorf("cannot found packet handler, packetId %d", packetId)
	}

	//TODO:packet handler를 task로 집어넣기
	//session 생명 주기는 어케할 것인가
	task.New(func(i ...interface{}) (error, error) {
		handler.Run(session, packet)
		return nil, nil
	})

	return nil
}

func (h *TestSession) OnSend(session session.ISession, sentBytes []byte) error {
	return nil
}

func (h *TestSession) OnDisconnect(session session.ISession) error {
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
	builder.Handler(&TestSession{})

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
