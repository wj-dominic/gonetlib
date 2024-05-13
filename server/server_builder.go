package server

import (
	"github.com/wj-dominic/gonetlib/logger"
	"github.com/wj-dominic/gonetlib/util/network"
)

type ServerBuilder interface {
	Configuration(ServerInfo) ServerBuilder
	Logger(logger.Logger) ServerBuilder
	Handler(ServerHandler) ServerBuilder
	Build() Server
}

type serverBuilder struct {
	config  ServerInfo
	handler ServerHandler
	logger  logger.Logger
}

func NewServerBuilder() ServerBuilder {
	config := logger.NewLoggerConfig().
		WriteToConsole().
		MinimumLevel(logger.DebugLevel)
	defaultLogger := config.CreateLogger()

	return &serverBuilder{
		config: ServerInfo{
			Address:    network.Endpoint{IP: "0.0.0.0", Port: 50000},
			MaxSession: 100,
			Protocols:  network.TCP,
		},
		handler: newDefaultServerHandler(),
		logger:  defaultLogger,
	}
}

func (builder *serverBuilder) Configuration(config ServerInfo) ServerBuilder {
	builder.config.Address = config.Address
	builder.config.MaxSession = config.MaxSession
	builder.config.Protocols = config.Protocols
	return builder
}

func (builder *serverBuilder) Logger(logger logger.Logger) ServerBuilder {
	builder.logger = logger
	return builder
}

func (builder *serverBuilder) Handler(handler ServerHandler) ServerBuilder {
	builder.handler = handler
	return builder
}

func (builder *serverBuilder) Build() Server {
	server := newServer(builder.logger, builder.config, builder.handler)
	return server
}
