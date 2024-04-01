package mmo_server

import (
	"gonetlib/logger"
	"gonetlib/server"
	"gonetlib/util/network"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := logger.CreateLoggerConfig().
		WriteToConsole().
		WriteToFile(
			logger.WriteToFile{
				Filepath:        "./MMOServer.log",
				RollingInterval: logger.RollingIntervalDay,
			}).
		MinimumLevel(logger.DebugLevel).
		TickDuration(1000)
	_logger := config.CreateLogger()

	builder := server.CreateServerBuilder()
	builder.Configuration(server.ServerInfo{
		Id:         1,
		Address:    network.Endpoint{IP: "0.0.0.0", Port: 50000},
		Protocols:  network.TCP | network.UDP,
		MaxSession: 10000,
	})
	builder.Logger(_logger)
	builder.Handler(CreateMMOServer())

	server := builder.Build()
	server.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c

	server.Stop()
	_logger.Info("Success to stop the server", logger.Why("signal", sig.String()))
}
