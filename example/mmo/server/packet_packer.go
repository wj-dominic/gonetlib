package mmo_server

import "gonetlib/message"

type IPacketPacker interface {
	Unpack(*message.Message)
	Pack() *message.Message
	SetData(interface{})
	GetData() interface{}
}

type PacketPacker[T any] struct {
	structure T
}

func CreatePacker[T any]() IPacketPacker {
	return &PacketPacker[T]{}
}

func (p *PacketPacker[T]) Unpack(packet *message.Message) {
	packet.Pop(&p.structure)
}

func (p *PacketPacker[T]) Pack() *message.Message {
	packet := message.NewMessage()
	packet.Push(p.structure)
	return packet
}

func (p *PacketPacker[T]) SetData(structure interface{}) {
	p.structure = structure.(T)
}

func (p *PacketPacker[T]) GetData() interface{} {
	return p.structure
}
