package session

import (
	"bytes"
	"context"
	"fmt"
	"gonetlib/logger"
	"gonetlib/message"
	"gonetlib/ringbuffer"
	"time"
)

const (
	TickDurationForSend time.Duration = time.Second * 5
	DeadlineTimeForSend time.Duration = time.Second * 5
)

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
	defer func() {
		if session.release() == true {
			session.onRelease()
		}
	}()

	session.acquire()
	if session.conn == nil {
		return fmt.Errorf("connection is nil")
	}

	if session.handler != nil {
		if err := session.handler.OnConnect(); err != nil {
			session.logger.Error("Failed to connect ")
		}
	}

	session.wg.Add(2)
	go session.readAsync()
	go session.sendAsync()

	return nil
}

func (session *TCPSession) readAsync() {
	defer func() {
		session.wg.Done()
		if session.release() == true {
			session.onRelease()
		}
	}()

	session.acquire()
	for {
		//빈 버퍼 획득
		buffer := session.recvBuffer.GetRearBuffer()
		if buffer == nil {
			session.logger.Error("Failed to get for read buffer")
			return
		}

		recvBytes, err := session.conn.Read(buffer)
		if err != nil {
			return
		}

		if session.recvBuffer.MoveRear(uint32(recvBytes)) == false {
			session.logger.Error("Failed to move receive buffer", logger.Why("recvBytes", recvBytes))
			return
		}

		packet := message.NewMessage()

		for {
			if session.recvBuffer.GetUseSize() <= uint32(packet.GetHeaderSize()) {
				break
			}

			//net 헤더 사이즈만큼 Peek
			session.recvBuffer.Peek(packet.GetHeaderBuffer(), uint32(packet.GetHeaderSize()))

			//링버퍼에 패킷 사이즈만큼 없을 경우 핸들링 처리 안함
			expectedPacketSize := uint32(packet.GetHeaderSize() + packet.GetExpectedPayloadSize())
			if session.recvBuffer.GetUseSize() < expectedPacketSize {
				break
			}

			//패킷 사이즈만큼 있으므로 앞에서 Peek한 만큼 링버퍼 소모
			session.recvBuffer.MoveFront(uint32(packet.GetHeaderSize()))

			//헤더에 있는 Payload 사이즈만큼 데이터 복사
			session.recvBuffer.Read(packet.GetPayloadBuffer(), uint32(packet.GetExpectedPayloadSize()))
			packet.MoveRear(packet.GetExpectedPayloadSize())

			if session.handler != nil {
				if err := session.handler.OnRecv(packet); err != nil {
					session.logger.Error("Failed to call on receieve", logger.Why("error", err))
					return
				}
			}

			packet.Reset()
		}
	}
}

func (session *TCPSession) sendAsync() {
	defer func() {
		session.wg.Done()
		if session.release() == true {
			session.onRelease()
		}
	}()

	session.acquire()
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
	defer func() {
		if session.release() == true {
			session.onRelease()
		}
	}()

	session.acquire()
	session.conn.Close()
	return nil
}

func (session *TCPSession) onRelease() {
	close(session.sendChannel)
	session.recvBuffer.Clear()

	if session.event != nil {
		session.event.OnRelease(session)
	}

	session.wg.Wait()
	session.reset()
}
