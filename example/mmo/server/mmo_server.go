package mmo_server

import (
	"fmt"
	"gonetlib/logger"
	"gonetlib/message"
	"gonetlib/session"
	"gonetlib/util"
	"reflect"
	"sync"
	"time"
)

type MMOServer struct {
	logger logger.ILogger
	count  int32

	nodes sync.Map
}

func CreateMMOServer() *MMOServer {
	return &MMOServer{
		logger: nil,
		count:  0,
	}
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
	node := CreateNode(session)
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

func (s *MMOServer) OnRecv(session session.ISession, packet *message.Message) error {
	node, err := s.getNode(session.GetID())
	if err != nil {
		s.logger.Error("Failed to get node from on recv", logger.Why("error", err.Error()))
		return err
	}

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

	_packer := GetPacker(packetId)
	_packer.Unpack(packet)

	ctx := CreateContext(node, _packer.GetData())
	node.SetContext(ctx)

	handler := GetPacketHandler(packetId)
	if handler == nil {
		return fmt.Errorf("cannot found packet handler, packetId %d", packetId)
	}

	//0번 고루틴에서 핸들러를 동작하게 한다.
	ctx.Async(func(i ...interface{}) (interface{}, error) {
		handler := i[0].(IPacketHandler)
		handler(ctx)
		return nil, nil
	}, 0).Start(handler)

	return nil
}

func (s *MMOServer) OnSend(session session.ISession, sentBytes []byte) error {
	return nil
}

func (s *MMOServer) OnDisconnect(session session.ISession) error {
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

type RequestLogin struct {
	id   uint64
	name string
}

type ResponseLogin struct {
	result  uint8
	message string
}

func RequestLoginHandler(ctx IPacketContext) {
	//패킷 데이터 참조하기
	request := ctx.GetPacket().(RequestLogin)
	fmt.Println(request.id)
	fmt.Println(request.name)

	//비동기 작업 요청하기
	ctx.Async(func(i ...interface{}) (interface{}, error) {
		//실제로 수행되는 비동기 작업
		time.Sleep(time.Second * 5)
		return i[0].(int) + i[1].(int), nil
	}).Await(func(result interface{}, err error) {
		//비동기 끝나고 호출되는 콜백
		fmt.Println(result)

		response := ResponseLogin{
			result:  0,
			message: "hello world!",
		}

		//노드 획득 후 사용하기
		ctx.GetNode().Send(response)
	}).Start(1, 2)
}
