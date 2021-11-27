package netserver

import "gonetlib/logger"

type NetServer struct {
	logger *logger.Logger
}

func NewNetServer() *NetServer {
	return &NetServer{logger: nil}
}

func (s *NetServer) Start() error {
	s.logger = logger.GetLogger()
	// if err != nil {
	// 	return err
	// }

	s.logger.SetLogConfig(logger.Max, "", "")
	s.logger.Start()
	s.logger.Error("test")
	return nil
}

func (s *NetServer) Stop() error {
	s.logger.Stop()
	return nil
}
