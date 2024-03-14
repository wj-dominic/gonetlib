package netserver

import (
	"gonetlib/netlogger"
	"log"
	"net"
)

type NetServer struct {
	addr string

	sessionMgr *SessionManager
}

func Run(address string) error {

	// 1. start logger
	netlogger.SetFileName("./netserver")

	// 2. create NetServer
	server := &NetServer{
		addr:       address,
		sessionMgr: nil,
	}

	// 3. create SessionManager
	server.sessionMgr = NewSessionManager()
	server.sessionMgr.Run()

	tcpAddr, err := net.ResolveTCPAddr("tcp", server.addr)
	if err != nil {
		netlogger.Error(err.Error())
		log.Fatal(err.Error())
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		netlogger.Error(err.Error())
		log.Fatal(err.Error())
	}
	defer listener.Close()

	netlogger.Info("Server is running")

	for {
		conn, err := listener.Accept()
		if err != nil {
			netlogger.Error(err.Error())
			break
		}

		err = server.sessionMgr.RequestNewSession(conn)
		if err != nil {
			netlogger.Error(err.Error())
			conn.Close()

			continue
		}
	}

	return nil
}

func (s *NetServer) Stop() error {
	s.sessionMgr.Stop()
	return nil
}
