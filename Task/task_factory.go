package task

import (
	"gonetlib/message"
	"gonetlib/netlogger"
	"gonetlib/util/singleton"
	"reflect"
	"sync"
)

type TaskFactory struct {
	TaskRegisters sync.Map
}

func Factory() *TaskFactory {
	return singleton.GetInstance[TaskFactory]()
}

func (factory *TaskFactory) Init() {
	factory.TaskRegisters = sync.Map{}
}

func (factory *TaskFactory) CreateTask(taskID uint16, packet *message.Message) *ITask {
	taskRegister, exist := factory.TaskRegisters.Load(taskID)
	if exist == false {
		netlogger.GetLogger().Error("Not found task register | taskID[%d]", taskID)
		return nil
	}

	newTask := taskRegister.(ITaskRegister).CreateTask(packet)
	if newTask == nil {
		netlogger.GetLogger().Error("Failed to regist a task | taskID[%d] register[%s]", taskID, reflect.TypeOf(taskRegister).Name())
		return nil
	}

	return newTask
}

func (factory *TaskFactory) AddRegister(taskID uint16, register ITaskRegister) bool {
	if register == nil {
		netlogger.GetLogger().Error("Invalid register | taskID[%d] register[%s]", taskID, reflect.TypeOf(register).Name())
		return false
	}

	if _, exist := factory.TaskRegisters.Load(taskID); exist == true {
		netlogger.GetLogger().Warn("Already has register in the factory | taskID[%d] register[%s]", taskID, reflect.TypeOf(register).Name())
		return false
	}

	factory.TaskRegisters.Store(taskID, register)

	return true
}
