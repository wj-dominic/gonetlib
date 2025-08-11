package client

import (
	"github.com/wj-dominic/gonetlib/logger"
	"github.com/wj-dominic/gonetlib/message"
	"github.com/wj-dominic/gonetlib/session"
)

type ClientHandler interface {
	OnRun(logger.Logger) error
	OnStop() error
	session.SessionHandler
}

type defaultClientHandler struct {
	logger logger.Logger
}

func newDefaultClientHandler() ClientHandler {
	return &defaultClientHandler{}
}

func (handler *defaultClientHandler) OnRun(logger logger.Logger) error {
	handler.logger = logger
	handler.logger.Warn("Client has been started by the default client handler. need to assign your client handler")
	return nil
}

func (handler *defaultClientHandler) OnStop() error {
	handler.logger.Warn("Client has been stopped by the default client handler. need to assign your client handler")
	return nil
}

func (handler *defaultClientHandler) OnConnect(session session.Session) error {
	handler.logger.Warn("Client has been connected to the session by the default client handler. need to assign your client handler", logger.Why("id", session.GetID()))
	return nil
}

func (handler *defaultClientHandler) OnRecv(session session.Session, packet *message.Message) error {
	handler.logger.Warn("Data has been received from the session by the default client handler. need to assign your client handler", logger.Why("id", session.GetID()))
	return nil
}

func (handler *defaultClientHandler) OnSend(session session.Session, sentBytes []byte) error {
	handler.logger.Warn("The default client handler has transmitted data to the session. need to assign your client handler", logger.Why("id", session.GetID()))
	return nil
}

func (handler *defaultClientHandler) OnDisconnect(session session.Session) error {
	handler.logger.Warn("Client has been connected from the session by the default client handler. need to assign your client handler", logger.Why("id", session.GetID()))
	return nil
}
