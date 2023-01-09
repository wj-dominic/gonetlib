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
