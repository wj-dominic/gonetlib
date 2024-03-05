package server

import "fmt"

type Endpoint struct {
	IP   string
	Port uint16
}

func (e *Endpoint) ToString() string {
	return fmt.Sprintf("%s:%d", e.IP, e.Port)
}
