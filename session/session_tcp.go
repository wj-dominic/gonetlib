package session

import (
	"bytes"
	"context"
	"fmt"
	"gonetlib/logger"
	"gonetlib/ringbuffer"
	"time"
)

const (
	TickDurationForSend time.Duration = time.Second * 5
	DeadlineTimeForSend time.Duration = time.Second * 5
)

type ITCPSessionHandler interface {
	OnRelease(id uint64)
}

type TCPSession struct {
	Session
	recvBuffer  *ringbuffer.RingBuffer
	sendChannel chan []byte
}

func newTcpSession(logger logger.ILogger, ctx context.Context) ISession {
	return &TCPSession{
		Session:     newSession(logger, ctx),
		recvBuffer:  ringbuffer.NewRingBuffer(true),
		sendChannel: make(chan []byte),
	}
}

func (session *TCPSession) Start() error {
	if session.conn == nil {
		return fmt.Errorf("connection is nil")
	}

	if session.handler != nil {
		return session.handler.OnConnect()
	}

	session.wg.Add(2)
	go session.readAsync()
	go session.sendAsync()

	return nil
}

func (session *TCPSession) readAsync() {
	defer func() {
		session.wg.Done()
		//session.release()
	}()

	for {
		//빈 버퍼 획득
		buffer := session.recvBuffer.GetRearBuffer()

		_, err := session.conn.Read(buffer)
		if err != nil {
			return
		}
	}
}

func (session *TCPSession) sendAsync() {
	defer session.wg.Done()

	ontick := time.NewTicker(TickDurationForSend)
	var sendBuffer bytes.Buffer
Loop:
	for {
		select {
		case <-session.ctx.Done():
			break Loop

		case msg := <-session.sendChannel:
			sendBuffer.Write(msg)

		case <-ontick.C:
			sentBytes, err := session.sendBytes(sendBuffer.Bytes())
			if err != nil {
				session.logger.Error("Failed to send",
					logger.Why("error", err.Error()))
				break Loop
			}

			if sendBuffer.Len() != sentBytes {
				session.logger.Error("Invalid bytes for sent",
					logger.Why("sendBytes", sendBuffer.Len()),
					logger.Why("sentBytes", sentBytes))
				break Loop
			} else {
				sendBuffer.Reset()
			}
		}
	}

	session.logger.Debug("End async send", logger.Why("id", session.GetID()))
}

func (session *Session) sendBytes(data []byte) (int, error) {
	if len(data) <= 0 {
		return 0, fmt.Errorf("empty data for send")
	}

	if err := session.conn.SetWriteDeadline(time.Now().Add(DeadlineTimeForSend)); err != nil {
		return 0, err
	}

	sendBytes, err := session.conn.Write(data)
	if err != nil {
		return sendBytes, err
	}

	return sendBytes, nil
}

func (session *TCPSession) Stop() error {
	return nil
}
