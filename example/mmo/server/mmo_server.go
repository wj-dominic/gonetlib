package mmo_server

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/wj-dominic/gonetlib/logger"
	"github.com/wj-dominic/gonetlib/message"
	"github.com/wj-dominic/gonetlib/session"
	"github.com/wj-dominic/gonetlib/util"
)

type MMOServer struct {
	logger logger.Logger
	count  int32

	nodes sync.Map
}

func NewMMOServer() *MMOServer {
	return &MMOServer{
		logger: nil,
		count:  0,
	}
}

func (s *MMOServer) OnRun(logger logger.Logger) error {
	s.logger = logger
	s.count = 0
	return nil
}

func (s *MMOServer) OnStop() error {
	s.logger.Info("Connected session count", logger.Why("count", s.count))
	return nil
}

func (s *MMOServer) OnConnect(session session.Session) error {
	node := NewNode(session)
	s.addNode(session.GetID(), node)

	s.logger.Info("On connect session", logger.Why("id", session.GetID()))
	util.InterlockIncrement(&s.count)
	return nil
}

func (s *MMOServer) addNode(sessionID uint64, node INode) error {
	if _, loaded := s.nodes.LoadOrStore(sessionID, node); loaded == true {
		return fmt.Errorf("already has same node in container | sessionID[%d]", sessionID)
	}

	return nil
}

func (s *MMOServer) getNode(sessionID uint64) (INode, error) {
	node, exist := s.nodes.Load(sessionID)
	if exist == false {
		return nil, fmt.Errorf("there is no node in container | sessionID[%d]", sessionID)
	}

	return node.(INode), nil
}

func (s *MMOServer) removeNode(sessionID uint64) error {
	if _, err := s.getNode(sessionID); err != nil {
		return err
	}

	s.nodes.Delete(sessionID)
	return nil
}

func (s *MMOServer) OnRecv(session session.Session, packet *message.Message) error {
	node, err := s.getNode(session.GetID())
	if err != nil {
		s.logger.Error("Failed to get node from on recv", logger.Why("error", err.Error()))
		return err
	}

	header := PacketHeader{}
	headerSize := uint16(util.Sizeof(reflect.ValueOf(header)))

	packet.Peek(&header)

	if packet.GetPayloadSize() < headerSize+header.Length {
		return fmt.Errorf("payload size grater then packet size, payloadSize %d > packetSize %d", headerSize+header.Length, packet.GetPayloadSize())
	}

	packet.MoveFront(headerSize)

	if packet.GetPayloadSize() != header.Length {
		return fmt.Errorf("not matched payload size with packet size, payloadSize %d > packetSize %d", header.Length, packet.GetPayloadSize())
	}

	ctx, err := GetPacketContext(header.Id)
	if err != nil {
		s.logger.Error("Failed to get packet context", logger.Why("packetId", header.Id), logger.Why("sessionId", session.GetID()))
		return err
	}

	if err = ctx.UnPack(packet); err != nil {
		s.logger.Error("Failed to unpack from packet", logger.Why("packetId", header.Id), logger.Why("sessionId", session.GetID()))
		return err
	}

	ctx.SetNode(node)
	node.SetContext(ctx)

	ctx.RunHandler(0)

	return nil
}

func (s *MMOServer) OnSend(session session.Session, sentBytes []byte) error {
	return nil
}

func (s *MMOServer) OnDisconnect(session session.Session) error {
	node, err := s.getNode(session.GetID())
	if err != nil {
		s.logger.Error("Failed to get node from on disconnect", logger.Why("error", err.Error()))
		return err
	}

	//task가 끝날 때까지 대기하기
	//이미 종료하는 상황이기 때문에 해당 세션을 담당하는 고루틴은 대기해도 된다.
	node.Wait()

	node.Clear()
	s.removeNode(session.GetID())

	s.logger.Info("On disconnect session", logger.Why("id", session.GetID()))
	return nil
}
