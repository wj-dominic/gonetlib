package gym

import "sync"

var gymManager *GymManager = nil
var once sync.Once

type GymManager struct {
	gyms	map[GymType]*Gym
}

func newGyms() *GymManager {
	return &GymManager{
		gyms : make(map[GymType]*Gym),
	}
}

func GetInstance() *GymManager {
	once.Do(func() {
		if gymManager == nil {
			gymManager = newGyms()
		}
	})

	return gymManager
}

func (gymManager *GymManager) CreateGym(gymType GymType, trainerCount uint8, routineCount uint8) bool {
	_, exist := gymManager.gyms[gymType]
	if exist == true {
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

