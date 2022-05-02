package singleton

import (
	"reflect"
	"sync"
)

type Instance interface{
	Init()
}

var cache 	sync.Map

func GetInstance[T any]() (t *T){
	typeName := reflect.TypeOf(t)
	if value, exist := cache.Load(typeName) ; exist == true{
		return value.(*T)
	}

	newValue := new(T)
	getValue, _ := cache.LoadOrStore(typeName, newValue)

	if _, ok := typeName.MethodByName("Init"); ok == true{
		instance := getValue.(Instance)
		instance.Init()
	}

	return getValue.(*T)
}
