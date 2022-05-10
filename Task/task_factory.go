package task

import (
	"gonetlib/message"
	"gonetlib/netlogger"
	"gonetlib/util/singleton"
	"reflect"
	"sync"
)

func CreateTask(taskID uint16, packet *message.Message) ITask {
	var factory *TaskFactory = singleton.GetInstance[TaskFactory]()
	taskRegister, exist := factory.taskRegisters.Load(taskID)
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

type TaskFactory struct {
	taskRegisters sync.Map
}

func (factory *TaskFactory) Init() {
	factory.taskRegisters = sync.Map{}
}

func (factory *TaskFactory) AddRegister(taskID uint16, register ITaskRegister) bool {
	if register == nil {
		netlogger.GetLogger().Error("Invalid register | taskID[%d] register[%s]", taskID, reflect.TypeOf(register).Name())
		return false
	}

	if _, exist := factory.taskRegisters.Load(taskID); exist == true {
		netlogger.GetLogger().Warn("Already has register in the factory | taskID[%d] register[%s]", taskID, reflect.TypeOf(register).Name())
		return false
	}

	factory.taskRegisters.Store(taskID, register)

	return true
}
