package netserver

import (
	. "gonetlib/netlogger"
	"log"
	"net"
)

type NetServer struct {
	addr string

	sessionMgr *SessionManager
}

func Run(address string) error {

	// 1. start logger
	GetLogger().SetLogConfig(Max, "", "")
	err := GetLogger().Start()
	if err != nil {
		log.Print(err.Error())
	}

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
		GetLogger().Error(err.Error())
		log.Fatal(err.Error())
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		GetLogger().Error(err.Error())
		log.Fatal(err.Error())
	}
	defer listener.Close()

	GetLogger().Info("Server is running")

	for {
		conn, err := listener.Accept()
		if err != nil {
			GetLogger().Error(err.Error())
			break
		}

		err = server.sessionMgr.RequestNewSession(conn)
		if err != nil {
			GetLogger().Error(err.Error())
			conn.Close()

			continue
		}
	}

	return nil
}

func (s *NetServer) Stop() error {
	s.sessionMgr.Stop()
	GetLogger().Stop()
	return nil
}
