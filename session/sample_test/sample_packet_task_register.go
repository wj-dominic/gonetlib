package sample_test

import (
	"gonetlib/message"
	"gonetlib/netlogger"
	"gonetlib/task"
)

//TODO :: 제너레이팅
type TASK_REGISTER_REQ_ECHO struct {
}

func (r *TASK_REGISTER_REQ_ECHO) CreateTask(packet *message.Message) task.ITask{
	if packet == nil {
		netlogger.GetLogger().Error("packet is nullptr")
		return nil
	}

	var newTask TASK_REQ_ECHO

	packet.Pop(&newTask)

	return &newTask
}

func AddTaskRegister_REQ_ECHO(){
	task.AddTaskRegister(REQ_ECHO, &TASK_REGISTER_REQ_ECHO{})
}