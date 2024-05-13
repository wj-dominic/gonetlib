package mmo_server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/wj-dominic/gonetlib/server"
	"github.com/wj-dominic/gonetlib/util/network"

	"github.com/wj-dominic/gonetlib/logger"
)

func main() {
	config := logger.NewLoggerConfig().
		WriteToConsole().
		WriteToFile(
			logger.WriteToFile{
				Filepath:        "./MMOServer.log",
				RollingInterval: logger.RollingIntervalDay,
			}).
		MinimumLevel(logger.DebugLevel).
		TickDuration(1000)
	_logger := config.CreateLogger()

	builder := server.NewServerBuilder()
	builder.Configuration(server.ServerInfo{
		Id:         1,
		Address:    network.Endpoint{IP: "0.0.0.0", Port: 50000},
		Protocols:  network.TCP | network.UDP,
		MaxSession: 10000,
	})
	builder.Logger(_logger)
	builder.Handler(NewMMOServer())

	server := builder.Build()
	server.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c

	server.Stop()
	_logger.Info("Success to stop the server", logger.Why("signal", sig.String()))
}
