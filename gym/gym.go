package gym

import (
	. "gonetlib/netlogger"
	. "gonetlib/trainer"
)

type GymType uint8
const (
	Main GymType = 0
)

type Gym struct{
	name       	string
	gymType		GymType

	routines 	[]chan string //TODO : routine 패키지 제작 필요
	trainers 	[]*Trainer     //TODO : routine 처리 스레드 제작 필요
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

	gym.routines = make([]chan string, routineCount)
	gym.trainers = make([]*Trainer, trainerCount)

	for index := range gym.trainers {
		id := uint8(index)
		gym.trainers[index] = NewTrainer(id)
	}

	return true
}

