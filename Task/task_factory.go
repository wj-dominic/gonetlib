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
	taskRegister := factory.GetTaskRegister(taskID)
	if taskRegister == nil{
		netlogger.GetLogger().Error("Failed to create task | has no [%d] register", taskID)
		return nil
	}

	newTask := taskRegister.CreateTask(packet)
	if newTask == nil {
		netlogger.GetLogger().Error("Failed to regist a task | taskID[%d] register[%s]", taskID, reflect.TypeOf(taskRegister).Name())
		return nil
	}

	return newTask
}

func AddTaskRegister(taskID uint16, register ITaskRegister) bool{
	var factory *TaskFactory = singleton.GetInstance[TaskFactory]()
	return factory.AddTaskRegister(taskID, register)
}

type TaskFactory struct {
	taskRegisters sync.Map
}

func (factory *TaskFactory) Init() {
	factory.taskRegisters = sync.Map{}
}

func (factory *TaskFactory) GetTaskRegister(taskID uint16) ITaskRegister{
	taskRegister, exist := factory.taskRegisters.Load(taskID)
	if exist == false {
		return nil
	}

	return taskRegister.(ITaskRegister)
}

func (factory *TaskFactory) AddTaskRegister(taskID uint16, register ITaskRegister) bool {
	if register == nil {
		netlogger.GetLogger().Error("Invalid register | taskID[%d] register[%s]", taskID, reflect.TypeOf(register).Name())
		return false
	}

	if factory.GetTaskRegister(taskID) != nil{
		netlogger.GetLogger().Warn("Already has register in the factory | taskID[%d] register[%s]", taskID, reflect.TypeOf(register).Name())
		return false
	}

	factory.taskRegisters.Store(taskID, register)

	return true
}
