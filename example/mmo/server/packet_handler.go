package mmo_server

import "gonetlib/message"

type IPacketHandler interface {
	Init(message.Message)
	Run(Node)
}

type PacketHandlers struct {
	handlers map[uint16]IPacketHandler
}

func CreatePacketHandlers() *PacketHandlers {
	return &PacketHandlers{
		handlers: make(map[uint16]IPacketHandler),
	}
}

func (h *PacketHandlers) AddPacketHandler(id uint16, handler IPacketHandler) {
	if _, exist := h.handlers[id]; exist == false {
		h.handlers[id] = handler
	}
}

func (h *PacketHandlers) GetPacketHandler(id uint16) IPacketHandler {
	if _, exist := h.handlers[id]; exist == false {
		return nil
	}

	return h.handlers[id]
}
