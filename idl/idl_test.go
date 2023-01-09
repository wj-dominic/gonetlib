package idl

import (
	"encoding/binary"
	"fmt"
	reflect "reflect"
	"testing"

	"google.golang.org/protobuf/proto"
)

type IPacket interface {
	GetID() uint32
	GetMessage() interface{}
}

type PacketCreator struct {
	concreteMap map[uint32]func() IPacket
}

func NewPacketCreator() *PacketCreator {
	return &PacketCreator{
		concreteMap: make(map[uint32]func() IPacket),
	}
}

func (c *PacketCreator) Add(id uint32, constructor func() IPacket) error {
	if _, exist := c.concreteMap[id]; exist {
		return fmt.Errorf("duplicate id | id[%d]", id)
	}

	c.concreteMap[id] = constructor
	return nil
}

func (c *PacketCreator) Create(id uint32) (IPacket, error) {
	constructor, exist := c.concreteMap[id]
	if !exist {
		return nil, fmt.Errorf("invalid id | id[%d]", id)
	}

	return constructor(), nil
}

type ISerializer interface {
	Serialize(packet IPacket) ([]byte, error)
	Deserialize(buf []byte) (IPacket, error)
}

type ProtobufSerializer struct {
	packetCreator *PacketCreator
}

func NewProtobufSerializer(c *PacketCreator) ISerializer {
	return &ProtobufSerializer{
		packetCreator: c,
	}
}

func (ps *ProtobufSerializer) Serialize(packet IPacket) ([]byte, error) {
	out, err := proto.Marshal(packet.GetMessage().(proto.Message))
	if err != nil {
		return nil, err
	}

	id := packet.GetID()
	sizeOfId := int(reflect.TypeOf(id).Size())

	buffer := make([]byte, len(out)+sizeOfId)
	binary.LittleEndian.PutUint32(buffer, id)
	copy(buffer[sizeOfId:], out)

	return buffer, nil
}

func (ps *ProtobufSerializer) Deserialize(buf []byte) (IPacket, error) {
	sizeOfId := int(reflect.TypeOf((*uint32)(nil)).Elem().Size())
	id := binary.LittleEndian.Uint32(buf[0:sizeOfId])

	packet, err := ps.packetCreator.Create(id)
	if err != nil {
		return nil, fmt.Errorf("invalid packet id | id[%d]", id)
	}

	err = proto.Unmarshal(buf[sizeOfId:], packet.GetMessage().(proto.Message))
	if err != nil {
		return nil, err
	}

	return packet, nil
}

type Packet[TMessage any] struct {
	id  uint32
	msg TMessage
}

func NewPacket[TMessage any](id uint32) *Packet[TMessage] {
	return &Packet[TMessage]{
		id: id,
	}
}

func (p *Packet[TMessage]) GetID() uint32           { return p.id }
func (p *Packet[TMessage]) GetMessage() interface{} { return &p.msg }
func (p *Packet[TMessage]) Message() *TMessage      { return p.GetMessage().(*TMessage) }

func TestIDL(t *testing.T) {
	creator := NewPacketCreator()

	creator.Add(uint32(ID_REQ_ECHO), func() IPacket { return NewPacket[ReqEcho](uint32(ID_REQ_ECHO)) })
	creator.Add(uint32(ID_RES_ECHO), func() IPacket { return NewPacket[ResEcho](uint32(ID_RES_ECHO)) })

	serializer := NewProtobufSerializer(creator)

	packet := NewPacket[ReqEcho](uint32(ID_REQ_ECHO))
	packet.Message().From = "john"
	packet.Message().Id = ID_REQ_ECHO
	packet.Message().Message = "test"

	fmt.Printf("reqEcho: %v\n", packet.Message())

	//시리얼라이즈 proto to bytes
	out, err := serializer.Serialize(packet)
	if err != nil {
		t.Failed()
		return
	}

	fmt.Printf("out: %v\n", out)

	//디시리얼라이즈 bytes to proto
	msg, err := serializer.Deserialize(out)
	if err != nil {
		t.Failed()
		return
	}

	fmt.Printf("outReqEcho: %v\n", msg.GetMessage())
}
