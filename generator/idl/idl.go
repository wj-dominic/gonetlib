//TODO : 버전을 기입하고 해당 버전을 바탕으로 gen 될 수 있도록 수정

package idl

type PACKET_REQ_ECHO struct {
	Message string
}

type PACKET_RES_ECHO struct {
	Message string
}

type PAKCET_REQ_LOGIN struct {
	Id       string
	Password string
	Email    string
}

type PACKET_RES_LOGIN struct {
	Result string
}
