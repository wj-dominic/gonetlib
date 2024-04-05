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
	info    ClientInfo
	handler IClientHandler
	logger  logger.ILogger
}

func NewClientBuilder() ClientBuilder {
	return &clientBuilder{
		info: ClientInfo{
			ServerAddress: network.Endpoint{IP: "127.0.0.1", Port: 50000},
			Protocols:     network.TCP,
			ConnectorInfo: DefaultConnectorInfo(),
		},
	}
}

func (builder *clientBuilder) Configuration(info ClientInfo) ClientBuilder {
	builder.info.ServerAddress = info.ServerAddress
	builder.info.Protocols = info.Protocols

	if info.reconnectDuration != 0 {
		builder.info.reconnectDuration = info.reconnectDuration
	}

	if info.reconnectLimit != 0 {
		builder.info.reconnectLimit = info.reconnectLimit
	}

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
	return newClient(builder.logger, builder.info, builder.handler)
}
