package server

import (
	"github.com/wj-dominic/gonetlib/logger"
	"github.com/wj-dominic/gonetlib/message"
	"github.com/wj-dominic/gonetlib/session"
)

type ServerHandler interface {
	OnRun(logger.Logger) error
	OnStop() error
	session.SessionHandler
}

type defaultServerHandler struct {
	logger logger.Logger
}

func newDefaultServerHandler() ServerHandler {
	return &defaultServerHandler{}
}

func (handler *defaultServerHandler) OnRun(logger logger.Logger) error {
	handler.logger = logger
	handler.logger.Warn("Server has been started by the default server handler. need to assign your server handler")
	return nil
}

func (handler *defaultServerHandler) OnStop() error {
	handler.logger.Warn("Server has been stopped by the default server handler. need to assign your server handler")
	return nil
}

func (handler *defaultServerHandler) OnConnect(session session.Session) error {
	handler.logger.Warn("Session has been connected to the default server handler. need to assign your server handler", logger.Why("id", session.GetID()))
	return nil
}

func (handler *defaultServerHandler) OnRecv(session session.Session, packet *message.Message) error {
	handler.logger.Warn("Data has been received from the session by the default server handler. need to assign your server handler", logger.Why("id", session.GetID()))
	return nil
}

func (handler *defaultServerHandler) OnSend(session session.Session, sentBytes []byte) error {
	handler.logger.Warn("The default server handler has transmitted data to the session. need to assign your server handler", logger.Why("id", session.GetID()))
	return nil
}

func (handler *defaultServerHandler) OnDisconnect(session session.Session) error {
	handler.logger.Warn("Session has been disconnected to the default server handler. need to assign your server handler", logger.Why("id", session.GetID()))
	return nil
}
