package node

import (
	. "gonetlib/gym"
	. "gonetlib/message"
	. "gonetlib/netlogger"
	. "gonetlib/routine"
	"gonetlib/util"
	"reflect"
)

type nodeHeader struct{
	packetID	uint32
	length		uint16
}

type Node struct {

}



func (node *Node) OnRecv(packet *Message) bool {
	var header nodeHeader

	//컨텐츠 파트의 헤더를 확인한다.
	//패킷에서 헤더 값을 뽑는다. 패킷에 헤더에 정의된 길이 만큼의 데이터가 없을 수 있으므로 Peek으로 뽑느다.
	packet.Peek(&header)

	headerSize := uint16(util.Sizeof(reflect.ValueOf(header)))
	payloadSize := headerSize + header.length

	if uint16(packet.GetPayloadLength()) < payloadSize {
		GetLogger().Error("invalid payload length | packet[%d] header[%d]", packet.GetPayloadLength(), payloadSize)
		return false
	}

	//헤더 길이만큼 뽑는다.
	packet.MoveFront(uint32(headerSize))

	if uint16(packet.GetPayloadLength()) != header.length {
		GetLogger().Error("packet length differ data and header | pop[%d] header[%d]", packet.GetPayloadLength(), header.length)
		return false
	}

	//packetID에 맞는 루틴을 생성한다.
	routine := GetRoutineMaker().MakeRoutine(header.packetID, packet)

	//루틴을 체육관에 넣는다.
	if GetGyms().Insert(GymMain, routine, 0) == false {
		GetLogger().Error("failed to insert a routine")
		return false
	}

	return true
}