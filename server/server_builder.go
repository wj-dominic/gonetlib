package server

import (
	"gonetlib/logger"
)

type IServerBuilder interface {
	Configuration(ServerInfo) IServerBuilder
	Logger(logger.ILogger) IServerBuilder
	Handler(IServerHandler) IServerBuilder
	Build() IServer
}

type serverBuilder struct {
	config  ServerInfo
	handler IServerHandler
	logger  logger.ILogger
}

func CreateServerBuilder() IServerBuilder {
	return &serverBuilder{
		config: ServerInfo{
			Address:    Endpoint{IP: "0.0.0.0", Port: 50000},
			MaxSession: 100,
			Protocols:  TCP,
		},
	}
}

func (builder *serverBuilder) Configuration(config ServerInfo) IServerBuilder {
	builder.config.Address = config.Address
	builder.config.MaxSession = config.MaxSession
	builder.config.Protocols = config.Protocols
	return builder
}

func (builder *serverBuilder) Logger(logger logger.ILogger) IServerBuilder {
	builder.logger = logger
	return builder
}

func (builder *serverBuilder) Handler(handler IServerHandler) IServerBuilder {
	builder.handler = handler
	return builder
}

func (builder *serverBuilder) Build() IServer {
	server := newServer(builder.logger, builder.config, builder.handler)
	return server
}
