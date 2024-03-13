package server

import (
	"context"
	"gonetlib/logger"
)

type IServerBuilder interface {
	Configuration(ServerConfig) IServerBuilder
	Logger(logger.ILogger) IServerBuilder
	Handler(IServerHandler) IServerBuilder
	Build() IServer
}

type serverBuilder struct {
	config  ServerConfig
	handler IServerHandler
	logger  logger.ILogger
}

func CreateServerBuilder() IServerBuilder {
	return &serverBuilder{
		config: ServerConfig{
			Address:    Endpoint{IP: "0.0.0.0", Port: 50000},
			MaxSession: 100,
			Protocols:  TCP,
		},
	}
}

func (builder *serverBuilder) Configuration(config ServerConfig) IServerBuilder {
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
	//context 생성
	server := newServerWithContext(builder.config, context.Background())
	return server
}
