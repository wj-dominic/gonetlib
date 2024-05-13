package server_test

import (
	"testing"
	"time"

	"github.com/wj-dominic/gonetlib/logger"
	"github.com/wj-dominic/gonetlib/server"
	"github.com/wj-dominic/gonetlib/util/network"
)

func TestAcceptor(t *testing.T) {
	config := logger.NewLoggerConfig().
		WriteToConsole().
		WriteToFile(
			logger.WriteToFile{
				Filepath:        "./acceptor.log",
				RollingInterval: logger.RollingIntervalDay,
			}).
		MinimumLevel(logger.DebugLevel).
		TickDuration(1000)
	_logger := config.CreateLogger()

	acceptor := server.NewAcceptor(_logger, network.TCP|network.UDP, network.Endpoint{IP: "0.0.0.0", Port: 50000}, nil)
	acceptor.Start()

	time.Sleep(time.Second * 10)

	acceptor.Stop()
}
