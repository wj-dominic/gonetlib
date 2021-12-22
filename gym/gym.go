package gym

import (
	. "gonetlib/netlogger"
	. "gonetlib/routine"
	. "gonetlib/trainer"
)

type GymType uint8
const (
	GymMain GymType = 0
)

const (
	maxRoutines uint16 = 300
)

type Gym struct{
	name       	string
	gymType		GymType

	routines 	[]chan Routine
	trainers 	[]*Trainer
}

func NewGym(gymName string, gymType GymType) *Gym {
	return &Gym{
		name : gymName,
		gymType: gymType,

		trainers: nil,
	}
}

func (gym *Gym) Create(trainerCount uint8 ,routinesCount uint8) bool {
	if routinesCount == 0 {
		GetLogger().Error("routine count is zero")
		return false
	}

	if trainerCount == 0 {
		GetLogger().Error("trainer count is zero")
		return false
	}

	gym.routines = make([]chan Routine, routinesCount)
	gym.trainers = make([]*Trainer, trainerCount)

	//루틴 채널 생성
	for index := range gym.routines{
		gym.routines[index] = make(chan Routine, maxRoutines)
	}

	//트레이너 생성
	for index := range gym.trainers {
		id := uint8(index)
		routineNumber := id % routinesCount

		gym.trainers[index] = NewTrainer(id, &gym.routines[routineNumber])

		if gym.trainers[index] == nil {
			GetLogger().Error("failed to create trainer | id[%d]", id)
			return false
		}

		gym.trainers[index].Start()
	}

	return true
}

func (gym *Gym) Insert(routine Routine, trainerID uint8) bool {
	if uint8(len(gym.routines)) < trainerID {
		GetLogger().Error("cannot found a trainer | trainerID[%d] trainerCount[%d]", trainerID, len(gym.routines))
		return false
	}

	gym.routines[trainerID] <- routine
	return true
}
