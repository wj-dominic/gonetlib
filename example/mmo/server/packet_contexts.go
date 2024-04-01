package mmo_server

import "fmt"

var contexts map[uint16]IPacketContext = make(map[uint16]IPacketContext)

func AddPacketContext(id uint16, ctx IPacketContext) error {
	if _, exist := contexts[id]; exist == true {
		return fmt.Errorf("already has packet context, id:%d", id)
	}

	contexts[id] = ctx
	return nil
}

func GetPacketContext(id uint16) (IPacketContext, error) {
	ctx, exist := contexts[id]
	if exist == false {
		return nil, fmt.Errorf("there is no packet context, id:%d", id)
	}

	return ctx, nil
}
