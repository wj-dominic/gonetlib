package netserver

import "net"

type Session struct{
	sessionID	uint64			//세션 ID
	Conn		net.Conn		//TCP connection
	RecvChannel chan *Message	//수신 버퍼
	SendChannel chan *Message	//송신 버퍼
}