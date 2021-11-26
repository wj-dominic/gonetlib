package singleton

import (
	"sync"
)

var singleton *Singleton = nil
var once sync.Once

func GetSingleton() *Singleton {
	once.Do(func(){
		if singleton == nil {
			singleton = new(Singleton)
		}
	})

	return singleton
}

type Singleton struct{
	instanceMap sync.Map
}

func (singleton *Singleton) SetInstance(name string, value interface{}) {
	singleton.instanceMap.Store(name, value)
}

func (singleton *Singleton) GetInstance(name string) interface{} {
	if value, exist := singleton.instanceMap.Load(name) ; exist == true {
		return value
	}

	return nil
}