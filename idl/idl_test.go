package idl

import (
	"encoding/binary"
	"fmt"
	"reflect"
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

func NewProtobufPacket(id uint32, message proto.Message) IPacket[proto.Message] {
	return &ProtobufPacket{
		ID:      id,
		Message: message,
	}
}

// GetID implements IPacket
func (p *ProtobufPacket) GetID() uint32 {
	return p.ID
}

// GetMessage implements IPacket
func (p *ProtobufPacket) GetMessage() proto.Message {
	return p.Message
}

type ProtobufPacketRegister struct {
	registerMap map[uint32]func() proto.Message
}

func NewProtobufPacketRegister() *ProtobufPacketRegister {
	return &ProtobufPacketRegister{
		registerMap: make(map[uint32]func() proto.Message),
	}
}

func (r *ProtobufPacketRegister) Regist(id uint32, constructor func() proto.Message) error {
	if _, exist := r.registerMap[id]; exist {
		return fmt.Errorf("duplicate id | id[%d]", id)
	}

	r.registerMap[id] = constructor

	return nil
}

func (r *ProtobufPacketRegister) GetPacket(id uint32) (proto.Message, error) {
	constructor, exist := r.registerMap[id]
	if !exist {
		return nil, fmt.Errorf("invalid id | id[%d]", id)
	}

	return constructor(), nil
}

type ISerializer[TInput any, TOutput any] interface {
	Serialize(data TInput) (TOutput, error)
	Deserialize(data TOutput) (TInput, error)
}

type ProtobufSerializer struct {
	packetRegister *ProtobufPacketRegister
}

func NewProtobufSerializer(register *ProtobufPacketRegister) ISerializer[IPacket[proto.Message], []byte] {
	return &ProtobufSerializer{
		packetRegister: register,
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

	msg, err := s.packetRegister.GetPacket(id)
	if err != nil {
		return nil, fmt.Errorf("invalid packet id | id[%d]", id)
	}

	err = proto.Unmarshal(data[sizeOfId:], msg)
	if err != nil {
		return nil, err
	}

	packet := NewProtobufPacket(id, msg)

	return packet, nil
}

func NewReqEcho() proto.Message { return &ReqEcho{} }
func NewResEcho() proto.Message { return &ResEcho{} }

func TestIDL(t *testing.T) {
	register := NewProtobufPacketRegister()
	register.Regist(uint32(ID_REQ_ECHO), NewReqEcho)
	register.Regist(uint32(ID_RES_ECHO), NewResEcho)

	serializer := NewProtobufSerializer(register)

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

	echo := msg.GetMessage().(*ReqEcho)

	fmt.Printf("outReqEcho: %s\n", echo.GetFrom())
}
