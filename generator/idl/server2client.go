package server2client

import "container/list"

type PACKET_REQ_MESSAGE struct {
	Message string
	Id      int
}

type PACKET_RES_MESSAGE struct {
	Message string
	Bulk    string
}

type PACKET_REQ_TEST struct {
	Id    int
	Names []string
	Map   map[string]int
	List  list.List
}
