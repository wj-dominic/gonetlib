package client

import (
	"gonetlib/logger"
	"gonetlib/util/network"
)

type IClientBuilder interface {
	Configuration(ClientInfo) IClientBuilder
	Logger(logger.ILogger) IClientBuilder
	Handler(IClientHandler) IClientBuilder
	Build() IClient
}

type clientBuilder struct {
	config  ClientInfo
	handler IClientHandler
	logger  logger.ILogger
}

func CreateClientBuilder() IClientBuilder {
	return &clientBuilder{
		config: ClientInfo{
			ServerAddress: network.Endpoint{IP: "127.0.0.1", Port: 50000},
			Protocols:     network.TCP,
		},
	}
}

func (builder *clientBuilder) Configuration(config ClientInfo) IClientBuilder {
	builder.config.ServerAddress = config.ServerAddress
	builder.config.Protocols = config.Protocols
	return builder
}

func (builder *clientBuilder) Logger(logger logger.ILogger) IClientBuilder {
	builder.logger = logger
	return builder
}

func (builder *clientBuilder) Handler(handler IClientHandler) IClientBuilder {
	builder.handler = handler
	return builder
}

func (builder *clientBuilder) Build() IClient {
	return newClient(builder.logger, builder.config, builder.handler)
}
