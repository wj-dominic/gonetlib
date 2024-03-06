package node

import (
	"gonetlib/message"
	"gonetlib/netlogger"
	"gonetlib/task"
	"gonetlib/util"
	"reflect"
)

type NodeHeader struct {
	PacketID uint16
	Length   uint16
}

type UserNode struct {
	session ISession
}

func NewUserNode(session ISession) *UserNode {
	if session == nil {
		return nil
	}

	return &UserNode{
		session: session,
	}
}

func (node *UserNode) OnConnect() {

}

func (node *UserNode) OnDisconnect() {

}

func (node *UserNode) OnSend(sendBytes int) {
	netlogger.Info("Send bytes | size[%d]", sendBytes)
}

func (node *UserNode) OnRecv(packet *message.Message) bool {
	var header NodeHeader

	//컨텐츠 파트의 헤더를 확인한다.
	//패킷에서 헤더 값을 뽑는다. 패킷에 헤더에 정의된 길이 만큼의 데이터가 없을 수 있으므로 Peek으로 뽑느다.
	packet.Peek(&header)

	headerSize := uint16(util.Sizeof(reflect.ValueOf(header)))
	payloadSize := headerSize + header.Length

	if uint16(packet.GetPayloadSize()) < payloadSize {
		netlogger.Error("invalid payload Length | packet[%d] header[%d]", packet.GetPayloadSize(), payloadSize)
		return false
	}

	//헤더 길이만큼 뽑는다.
	packet.MoveFront(headerSize)

	//헤더에 있는 길이와 패킷에 남아있는 데이터의 길이가 일치해야 한다.
	if uint16(packet.GetPayloadSize()) != header.Length {
		netlogger.Error("packet Length differ data and header | pop[%d] header[%d]", packet.GetPayloadSize(), header.Length)
		return false
	}

	newTask := task.CreateTask(header.PacketID, packet)
	bucket := task.GetBucket(0)
	if bucket == nil {
		netlogger.Error("Not found bucket | bucketID[%d]", 0)
		return false
	}

	if bucket.AddTask(newTask, 0) == false {
		netlogger.Error("failed to insert a task")
		return false
	}

	return true
}

func (node *UserNode) Send(header NodeHeader, value interface{}) bool {
	packet := message.NewMessage()

	packet.Push(header)
	packet.Push(value)

	packet.MakeHeader()

	return node.session.SendPost(packet)
}
