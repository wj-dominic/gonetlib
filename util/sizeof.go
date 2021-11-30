package util

import (
	"reflect"
	"sync"
)

var structSize sync.Map // map[reflect.Type]int

func Sizeof(value reflect.Value) int {
	switch value.Kind() {
	case reflect.Slice:
		if s := sizeof(value.Type().Elem()); s >= 0 {
			return s * value.Len()
		}
		return -1

	case reflect.String:
		return value.Len()

	case reflect.Struct:
		t := value.Type()
		if size, ok := structSize.Load(t); ok {
			return size.(int)
		}
		size := sizeof(t)
		structSize.Store(t, size)
		return size

	default:
		return sizeof(value.Type())
	}
}

func sizeof(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Array:
		if s := sizeof(t.Elem()); s >= 0 {
			return s * t.Len()
		}

	case reflect.Struct:
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			s := sizeof(t.Field(i).Type)
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return int(t.Size())
	}

	return -1
}