package client

import (
	"gonetlib/logger"
	"gonetlib/util/network"
)

type ClientBuilder interface {
	Configuration(ClientInfo) ClientBuilder
	Logger(logger.ILogger) ClientBuilder
	Handler(IClientHandler) ClientBuilder
	Build() Client
}

type clientBuilder struct {
	config  ClientInfo
	handler IClientHandler
	logger  logger.ILogger
}

func NewClientBuilder() ClientBuilder {
	return &clientBuilder{
		config: ClientInfo{
			ServerAddress: network.Endpoint{IP: "127.0.0.1", Port: 50000},
			Protocols:     network.TCP,
		},
	}
}

func (builder *clientBuilder) Configuration(config ClientInfo) ClientBuilder {
	builder.config.ServerAddress = config.ServerAddress
	builder.config.Protocols = config.Protocols
	return builder
}

func (builder *clientBuilder) Logger(logger logger.ILogger) ClientBuilder {
	builder.logger = logger
	return builder
}

func (builder *clientBuilder) Handler(handler IClientHandler) ClientBuilder {
	builder.handler = handler
	return builder
}

func (builder *clientBuilder) Build() Client {
	return newClient(builder.logger, builder.config, builder.handler)
}
