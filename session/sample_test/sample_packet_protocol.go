package sample_test

//TODO :: 사용자 정의 영역
type PACKET_REQ_ECHO struct {
	Message string
}

func NEW_PACKET_REQ_ECHO(_message string) (uint16, PACKET_REQ_ECHO) {
	return REQ_ECHO, PACKET_REQ_ECHO{
		Message: _message,
	}
}

type PACKET_RES_ECHO struct {
	Message string
}
