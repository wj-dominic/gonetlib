package routine

import (
	. "gonetlib/message"
	. "gonetlib/netlogger"
	. "gonetlib/singleton"
)

const (
	routineMakerName string = "ROUTINE_MAKER"
)

type routineRegister interface{
	Make(packet *Message) Routine
}

type routineMaker struct{
	routineFactory	map[uint32]routineRegister
}

func newMaker() {
	routineMaker := &routineMaker{
		routineFactory: make(map[uint32]routineRegister),
	}

	s := GetSingleton()
	s.SetInstance(routineMakerName, routineMaker)
}

func GetRoutineMaker() *routineMaker {
	s := GetSingleton()

	if s.GetInstance(routineMakerName) == nil {
		newMaker()
	}

	return s.GetInstance(routineMakerName).(*routineMaker)
}

// MakeRoutine : id에 해당 되는 루틴 생성 함수를 찾고 여기서 루틴을 생성한다.
func (maker *routineMaker) MakeRoutine(id uint32, packet *Message) Routine {
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
func (maker *routineMaker) AddRegister(id uint32, register routineRegister) bool {
	if _, exist := maker.routineFactory[id] ; exist == true {
		GetLogger().Error("already has routine register | id[%d]", id)
		return false
	}

	maker.routineFactory[id] = register

	return true
}