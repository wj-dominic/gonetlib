package client

import (
	"github.com/wj-dominic/gonetlib/logger"
	"github.com/wj-dominic/gonetlib/util/network"
)

type ClientBuilder interface {
	Configuration(ClientInfo) ClientBuilder
	Logger(logger.Logger) ClientBuilder
	Handler(ClientHandler) ClientBuilder
	Build() Client
}

type clientBuilder struct {
	info    ClientInfo
	handler ClientHandler
	logger  logger.Logger
}

func NewClientBuilder() ClientBuilder {
	config := logger.NewLoggerConfig().
		WriteToConsole().
		MinimumLevel(logger.DebugLevel)
	defaultLogger := config.CreateLogger()

	return &clientBuilder{
		info: ClientInfo{
			ServerAddress: network.Endpoint{IP: "127.0.0.1", Port: 50000},
			Protocols:     network.TCP,
			ConnectorInfo: DefaultConnectorInfo(),
		},
		handler: newDefaultClientHandler(),
		logger:  defaultLogger,
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

func (builder *clientBuilder) Logger(logger logger.Logger) ClientBuilder {
	builder.logger = logger
	return builder
}

func (builder *clientBuilder) Handler(handler ClientHandler) ClientBuilder {
	builder.handler = handler
	return builder
}

func (builder *clientBuilder) Build() Client {
	return newClient(builder.logger, builder.info, builder.handler)
}
