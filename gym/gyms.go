package gym

import (
	. "gonetlib/singleton"
)

const (
	GymsName string = "GYM"
)

type GymManager struct {
	gyms	map[GymType]*Gym
}

func newGyms() {
	gymManager := &GymManager{
		gyms : make(map[GymType]*Gym),
	}

	s := GetSingleton()
	s.SetInstance(GymsName, gymManager)
}

func GetGyms() *GymManager {
	s := GetSingleton()

	if s.GetInstance(GymsName) == nil {
		newGyms()
	}

	return s.GetInstance(GymsName).(*GymManager)
}

func (gymManager *GymManager) CreateGym(gymType GymType, trainerCount uint8, routineCount uint8) bool {
	_, exist := gymManager.gyms[gymType]
	if exist == true {
		//TODO : 로그
		return false
	}

	gymName := "GYM[" + string(gymType) + "]"

	gym := NewGym(gymName, gymType)

	if gym.Create(trainerCount, routineCount) == false {
		//TODO : 로그
		gym = nil
		return false
	}

	gymManager.gyms[gymType] = gym

	return true
}

