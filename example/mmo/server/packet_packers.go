package mmo_server

var packers map[uint16]IPacketPacker = make(map[uint16]IPacketPacker)

func GetPacker(packetId uint16) IPacketPacker {
	packer, exist := packers[packetId]
	if exist == false {
		return nil
	}

	return packer
}

func AddPacker(packetId uint16, packer IPacketPacker) {
	if _, exist := packers[packetId]; exist == true {
		return
	}

	packers[packetId] = packer
}
