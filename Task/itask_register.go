package task

import (
	"gonetlib/message"
)

type ITaskRegister interface {
	CreateTask(packet *message.Message) *ITask
}
