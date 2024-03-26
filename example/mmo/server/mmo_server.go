package mmo_server

import (
	"fmt"
	"gonetlib/logger"
	"gonetlib/message"
	"gonetlib/session"
	"gonetlib/task"
	"gonetlib/util"
	"reflect"
)

type MMOServer struct {
	logger   logger.ILogger
	handlers PacketHandlers
	count    int32
}

func (s *MMOServer) OnRun(logger logger.ILogger) error {
	s.logger = logger
	s.count = 0
	return nil
}

func (s *MMOServer) OnStop() error {
	s.logger.Info("Connected session count", logger.Why("count", s.count))
	return nil
}

func (s *MMOServer) OnConnect(session session.ISession) error {
	s.logger.Info("On connect session", logger.Why("id", session.GetID()))
	util.InterlockIncrement(&s.count)
	return nil
}

func (s *MMOServer) OnRecv(session session.ISession, packet *message.Message) error {
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
	handler := s.handlers.GetPacketHandler(packetId)
	if handler != nil {
		return fmt.Errorf("cannot found packet handler, packetId %d", packetId)
	}

	//session의 생명 주기와 다른 node를 둔다.
	task.New(func(i ...interface{}) (error, error) {
		//handler.Run(session)
		return nil, nil
	})

	return nil
}

func (s *MMOServer) OnSend(session session.ISession, sentBytes []byte) error {
	return nil
}

func (s *MMOServer) OnDisconnect(session session.ISession) error {
	s.logger.Info("On disconnect session", logger.Why("id", session.GetID()))
	return nil
}
