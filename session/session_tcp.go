package session

import (
	"bytes"
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
	gonetSession
	recvBuffer  *ringbuffer.RingBuffer
	sendChannel chan []byte
}

func NewTcpSession(logger logger.Logger) Session {
	return &TCPSession{
		gonetSession: newSession(logger),
		recvBuffer:   ringbuffer.NewRingBuffer(true),
		sendChannel:  make(chan []byte, 100),
	}
}

func (session *TCPSession) Start() error {
	defer func() {
		if session.release() == true {
			session.onRelease()
		}
	}()

	if session.acquire(true) == false {
		return fmt.Errorf("failed to acquire session for start, id[%d]", session.id)
	}

	session.logger.Debug("Start session", logger.Why("id", session.GetID()))

	if session.conn == nil {
		return fmt.Errorf("connection is nil")
	}

	session.wg.Add(2)
	if session.acquire() == true {
		go session.readAsync()
	}

	if session.acquire() == true {
		go session.sendAsync()
	}

	if session.handler != nil {
		if err := session.handler.OnConnect(session); err != nil {
			session.logger.Error("Failed to connect ")
		}
	}

	return nil
}

func (session *TCPSession) readAsync() {
	defer func() {
		session.wg.Done()
		if session.release() == true {
			session.onRelease()
		}
	}()

	session.logger.Debug("Begin read async", logger.Why("id", session.GetID()))

	for {
		//빈 버퍼 획득
		buffer := session.recvBuffer.GetRearBuffer()
		if buffer == nil {
			session.logger.Error("Failed to get for read buffer")
			break
		}

		recvBytes, err := session.conn.Read(buffer)
		if err != nil {
			session.logger.Error("Failed to read", logger.Why("error", err.Error()), logger.Why("id", session.GetID()))
			break
		}

		if session.recvBuffer.MoveRear(uint32(recvBytes)) == false {
			session.logger.Error("Failed to move receive buffer", logger.Why("recvBytes", recvBytes))
			break
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
				if err := session.handler.OnRecv(session, packet); err != nil {
					session.logger.Error("Failed to call on receieve", logger.Why("error", err))
					return
				}
			}

			session.monitoringData.RecvCount++
			session.monitoringData.RecvBytes += uint64(packet.GetExpectedPayloadSize())

			packet.Reset()
		}
	}

	session.logger.Debug("End read async", logger.Why("id", session.GetID()))
}

func (session *TCPSession) sendAsync() {
	defer func() {
		session.wg.Done()
		if session.release() == true {
			session.onRelease()
		}
	}()

	session.logger.Debug("Begin send async", logger.Why("id", session.GetID()))

	ontick := time.NewTicker(TickDurationForSend)
	var sendBuffer bytes.Buffer
Loop:
	for {
		select {
		case msg, ok := <-session.sendChannel:
			if ok == false {
				break Loop
			}

			sendBuffer.Write(msg)

		case <-ontick.C:
			if sendBuffer.Len() <= 0 {
				continue
			}

			sentBytes, err := session.sendBytes(sendBuffer.Bytes())
			if err != nil {
				session.logger.Error("Failed to send",
					logger.Why("error", err.Error()),
					logger.Why("sendBytes", sendBuffer.Len()))
				break Loop
			}

			if sendBuffer.Len() != sentBytes {
				session.logger.Error("Invalid bytes for sent",
					logger.Why("sendBytes", sendBuffer.Len()),
					logger.Why("sentBytes", sentBytes))
				break Loop
			} else {
				if session.handler != nil {
					session.handler.OnSend(session, sendBuffer.Bytes())
				}

				session.monitoringData.SendCount++
				session.monitoringData.SendBytes += uint64(sentBytes)

				sendBuffer.Reset()
			}
		}
	}

	session.logger.Debug("End send async", logger.Why("id", session.GetID()))
}

func (session *gonetSession) sendBytes(data []byte) (int, error) {
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
		//명시적인 종료이므로 stop을 호출한 스레드에서 대기
		session.wg.Wait()
		if session.release() == true {
			session.onRelease()
		}
	}()

	if session.acquire() == false {
		return fmt.Errorf("failed to acquire session for stop, id[%d]", session.id)
	}

	session.logger.Debug("Stop session", logger.Why("id", session.GetID()))

	if session.conn.Close() == nil {
		session.conn = nil
		if session.handler != nil {
			session.handler.OnDisconnect(session)
		}
	}

	if isClosed(session.sendChannel) == false {
		close(session.sendChannel)
	}

	return nil
}

func (session *TCPSession) Send(msg interface{}) {
	defer func() {
		if session.release() == true {
			session.onRelease()
		}
	}()

	if session.acquire() == false {
		return
	}

	packet := message.NewMessage()

	packet.Push(msg)

	packet.MakeHeader()

	session.sendChannel <- packet.GetBuffer()
}

func (session *TCPSession) SessionMonitoringData() SessionMonitoringData {
	return session.monitoringData
}

func (session *TCPSession) onRelease() {
	//한번만 실행되는 함수

	//1. 커넥션 종료
	if session.conn != nil {
		if session.conn.Close() == nil {
			session.conn = nil
			if session.handler != nil {
				session.handler.OnDisconnect(session)
			}
		}
	}

	//2. 샌드 채널 종료
	if isClosed(session.sendChannel) == false {
		close(session.sendChannel)
	}

	//3. 리시브 버퍼 정리
	session.recvBuffer.Clear()

	//4. 세션 ID 임시보관
	sessionID := session.GetID()

	//5. 고루틴 종료 대기 (수신, 송신 고루틴)
	session.wg.Wait()

	//6. 세션 초기화
	session.reset()

	session.logger.Debug("On release session", logger.Why("id", sessionID))

	//7. 전파
	if session.event != nil {
		session.event.OnRelease(sessionID, session)
	}
}

func isClosed(ch chan []byte) bool {
	notClosed := bool(true)
	select {
	case _, notClosed = <-ch:
	default:
	}

	return notClosed == false
}
