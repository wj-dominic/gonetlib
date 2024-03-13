package session

type TcpSession struct {
	Session
}

func newTcpSession() ISession {
	return &TcpSession{}
}

func (session *TcpSession) Start() {

}

func (session *TcpSession) Stop() {

}
