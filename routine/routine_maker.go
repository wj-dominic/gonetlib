package routine

import (
	. "gonetlib/message"
	. "gonetlib/netlogger"
	. "gonetlib/util"
)

type routineRegister interface{
	Make(packet *Message) Routine
}

type routineMaker struct{
	routineFactory	map[uint16]routineRegister
}

func (maker *routineMaker) Init(){
	maker.routineFactory = make(map[uint16]routineRegister)
}

func GetRoutineMaker() *routineMaker {
	return GetInstance[routineMaker]()
}

// MakeRoutine : id에 해당 되는 루틴 생성 함수를 찾고 여기서 루틴을 생성한다.
func (maker *routineMaker) MakeRoutine(id uint16, packet *Message) Routine {
	routineRegister, exist := maker.routineFactory[id]
	if exist == false {
		GetLogger().Error("failed to find a routine register | id[%d]", id)
		return nil
	}

	//루틴 생성 함수에서 make 한다.
	routine := routineRegister.Make(packet)
	if routine == nil {
		GetLogger().Error("failed to make routine | id[%d]", id)
		return nil
	}

	return routine
}

// AddRegister : 루틴마다 가지는 데이터가 다르므로 루틴 생성 함수를 등록한다. 이 함수는 서버가 시작될 때 여러 루틴을 등록해야 한다.
func (maker *routineMaker) AddRegister(id uint16, register routineRegister) bool {
	if _, exist := maker.routineFactory[id] ; exist == true {
		GetLogger().Error("already has routine register | id[%d]", id)
		return false
	}

	maker.routineFactory[id] = register

	return true
}