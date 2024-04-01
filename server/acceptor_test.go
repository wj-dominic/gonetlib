package server_test

import (
	"gonetlib/logger"
	"gonetlib/server"
	"gonetlib/util/network"
	"testing"
	"time"
)

func TestAcceptor(t *testing.T) {
	config := logger.CreateLoggerConfig().
		WriteToConsole().
		WriteToFile(
			logger.WriteToFile{
				Filepath:        "./acceptor.log",
				RollingInterval: logger.RollingIntervalDay,
			}).
		MinimumLevel(logger.DebugLevel).
		TickDuration(1000)
	_logger := config.CreateLogger()

	acceptor := server.CreateAcceptor(_logger, network.TCP|network.UDP, network.Endpoint{IP: "0.0.0.0", Port: 50000}, nil)
	acceptor.Start()

	time.Sleep(time.Second * 10)

	acceptor.Stop()
}
