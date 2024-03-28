package mmo_server

type IPacketHandler func(IPacketContext)

var handlers map[uint16]IPacketHandler = make(map[uint16]IPacketHandler)

func AddPacketHandler(id uint16, handler IPacketHandler) {
	if _, exist := handlers[id]; exist == false {
		handlers[id] = handler
	}
}

func GetPacketHandler(id uint16) IPacketHandler {
	if _, exist := handlers[id]; exist == false {
		return nil
	}

	return handlers[id]
}
