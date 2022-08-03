package idl

import (
	"encoding/binary"
	"fmt"
	reflect "reflect"
	"testing"

	"google.golang.org/protobuf/proto"
)

type IPacket[T any] interface {
	GetID() uint32
	GetMessage() T
}

type ProtobufPacket struct {
	ID      uint32
	Message proto.Message
}

// GetID implements IPacket
func (p *ProtobufPacket) GetID() uint32 {
	return p.ID
}

// GetMessage implements IPacket
func (p *ProtobufPacket) GetMessage() proto.Message {
	return p.Message
}

func NewProtobufPacket(id uint32, message proto.Message) IPacket[proto.Message] {
	return &ProtobufPacket{
		ID:      id,
		Message: message,
	}
}

type IPacketFactory[T any] interface {
	GetPacket(id uint32) T
}

type ProtobufPacketFactory struct {
}

// GetPacket implements IPacketFactory
func (f *ProtobufPacketFactory) GetPacket(id uint32) proto.Message {
	var msg proto.Message = nil

	switch ID(id) {
	case ID_REQ_ECHO:
		msg = new(ReqEcho)
		break
	case ID_RES_ECHO:
		msg = new(ResEcho)
		break
	default:
		break
	}

	return msg
}

func NewProtobufPacketFactory() IPacketFactory[proto.Message] {
	return &ProtobufPacketFactory{}
}

type ISerializer[TInput any, TOutput any] interface {
	Serialize(data TInput) (TOutput, error)
	Deserialize(data TOutput) (TInput, error)
}

type ProtobufSerializer struct {
	packetFactory IPacketFactory[proto.Message]
}

func NewProtobufSerializer() ISerializer[IPacket[proto.Message], []byte] {
	return &ProtobufSerializer{
		packetFactory: NewProtobufPacketFactory(),
	}
}

func (s *ProtobufSerializer) Serialize(data IPacket[proto.Message]) ([]byte, error) {
	out, err := proto.Marshal(data.GetMessage())
	if err != nil {
		return nil, err
	}

	id := data.GetID()
	sizeOfId := int(reflect.TypeOf(id).Size())

	buffer := make([]byte, len(out)+sizeOfId)
	binary.LittleEndian.PutUint32(buffer, id)
	copy(buffer[sizeOfId:], out)

	return buffer, nil
}

func (s *ProtobufSerializer) Deserialize(data []byte) (IPacket[proto.Message], error) {
	sizeOfId := int(reflect.TypeOf((*uint32)(nil)).Elem().Size())

	id := binary.LittleEndian.Uint32(data[0:sizeOfId])

	msg := s.packetFactory.GetPacket(id)
	if msg == nil {
		return nil, fmt.Errorf("invalid packet id | id[%d]", id)
	}

	err := proto.Unmarshal(data[sizeOfId:], msg)
	if err != nil {
		return nil, err
	}

	packet := NewProtobufPacket(id, msg)

	return packet, nil
}

func TestIDL(t *testing.T) {
	serializer := NewProtobufSerializer()

	reqEcho := &ReqEcho{From: "1", Message: "test"}
	fmt.Printf("reqEcho: %v\n", reqEcho)

	packet := NewProtobufPacket(uint32(ID_REQ_ECHO), reqEcho)

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

	echo := msg.GetMessage().ProtoReflect().Interface().(*ReqEcho)

	fmt.Printf("outReqEcho: %s\n", echo.GetFrom())
}
