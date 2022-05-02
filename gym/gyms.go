package gym

import (
	. "gonetlib/netlogger"
	. "gonetlib/routine"
	. "gonetlib/singleton"
)

const (
	GymsName string = "GYM"
)

type GymManager struct {
	gyms	map[GymType]*Gym
}

func newGyms() *GymManager{
	return &GymManager{
		gyms : make(map[GymType]*Gym),
	}
}

func GetGyms() *GymManager {
	return GetInstance[GymManager](newGyms)
}

func (gymManager *GymManager) CreateGym(gymType GymType, trainerCount uint8, routinesCount uint8) bool {
	_, exist := gymManager.gyms[gymType]
	if exist == true {
		GetLogger().Warn("already has gyms : " + string(gymType))
		return false
	}

	gymName := "GYM[" + string(gymType) + "]"

	gym := NewGym(gymName, gymType)

	if gym.Create(trainerCount, routinesCount) == false {
		GetLogger().Error("cannot create a gym")
		gym = nil
		return false
	}

	gymManager.gyms[gymType] = gym

	return true
}

func (gymManager *GymManager) Insert(gymType GymType, routine Routine, trainerID uint8) bool {
	if gymManager.gyms[gymType] == nil{
		GetLogger().Error("not found a gym | gymType[%d]", gymType)
		return false
	}

	if _, exist := gymManager.gyms[gymType] ; exist == false {
		GetLogger().Error("cannot found a gym | gymType[%d]", gymType)
		return false
	}

	return gymManager.gyms[gymType].Insert(routine, trainerID)
}
