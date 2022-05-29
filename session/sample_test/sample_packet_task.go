package sample_test

import "fmt"

//TODO :: 제너레이팅
type TASK_REQ_ECHO struct {
	PACKET_REQ_ECHO
}

func (t *TASK_REQ_ECHO) Run() bool {
	fmt.Println(t.Message)
	return true
}
