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

	nodes map[uint64]*Node
}

func CreateMMOServer() *MMOServer {
	return &MMOServer{
		logger:   nil,
		handlers: *CreatePacketHandlers(),
		count:    0,
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
	s.nodes[session.GetID()] = node

	//TODO:세션 접속했을 때 node 만들어서 session과 매핑하기
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
	if handler == nil {
		return fmt.Errorf("cannot found packet handler, packetId %d", packetId)
	}

	//TODO:시리얼라이저, context 만들기
	//[]byte => struct (deserialize)

	//struct => context
	ctx := CreateContext(node)

	//session의 생명 주기와 다른 node를 둔다.
	task.New(func(i ...interface{}) (error, error) {
		handler.Run(ctx)
		node.wait <- true
		//handler.Run(session)
		return nil, nil
	})

	return nil
}

func (s *MMOServer) OnSend(session session.ISession, sentBytes []byte) error {
	return nil
}

func (s *MMOServer) OnDisconnect(session session.ISession) error {
	node := s.nodes[session.GetID()]

	//waiting
	node.ctx.wait()

	node.session = nil

	s.logger.Info("On disconnect session", logger.Why("id", session.GetID()))
	return nil
}

// TODO:context에서 wait 어떻게 할 건지 고민 task wrapping??
type Context[TRequest any] struct {
	node    *Node
	request TRequest
	task    *Task
}

type RequestLogin struct {
	id   uint64
	name string
}

type LoginHandler struct {
}

func (h *LoginHandler) Run(ctx *Context[RequestLogin]) {

	//ctx.request.id
	//ctx.request.name

	//TODO:좀더 좋은 방법 고민
	ctx.task.async().await()

	node.wg.add(2)
	task.New(func(i ...interface{}) (string, error) {
		defer func() {
			node.wg.done()
		}()

		//DB 쿼리
		result := DB.Query("select...")
		return result
	}, 1).Await(func(result string, err2 error) {
		defer func() {
			node.wg.done()
		}()

		//DB 응답
		fmt.Println(result)
	})
}
