package netserver

import (
	"fmt"
)

type Address struct {
	IP   string
	Port uint16
}

func (Addr *Address) ToString() string {
	return fmt.Sprintf("%s:%d", Addr.IP, Addr.Port)
}

type NetServerSettings struct {
	Addr          Address
	MaxConnection uint64
}

func (setting *NetServerSettings) SetAddress(addr *Address) *NetServerSettings {
	setting.Addr = *addr
	return setting
}

func (setting *NetServerSettings) SetMaxConnection(max uint64) *NetServerSettings {
	setting.MaxConnection = max
	return setting
}

func (setting *NetServerSettings) Build() *NetServer {
	return &NetServer{}
}
