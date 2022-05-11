package sample_test


//TODO :: 제너레이팅
type PACKET_REQ_ECHO struct {
	message string
}

func NEW_PACKET_REQ_ECHO(_message string) (uint16, PACKET_REQ_ECHO){
	return REQ_ECHO, PACKET_REQ_ECHO{
		message: _message,
	}
}

type PACKET_RES_ECHO struct {
	message string
}


