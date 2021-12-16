package gym

import (
	. "gonetlib/netlogger"
	. "gonetlib/trainer"
	. "gonetlib/routine"
)

type GymType uint8
const (
	GymMain GymType = 0
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

func (gym *Gym) Create(routineCount uint8, trainerCount uint8) bool {
	if routineCount == 0 {
		GetLogger().Error("routine count is zero")
		return false
	}

	if trainerCount == 0 {
		GetLogger().Error("trainer count is zero")
		return false
	}

	gym.routines = make([]chan Routine, routineCount)
	gym.trainers = make([]*Trainer, trainerCount)

	for index := range gym.trainers {
		id := uint8(index)
		gym.trainers[index] = NewTrainer(id)
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
