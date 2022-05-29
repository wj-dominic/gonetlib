package util

import (
	"reflect"
)

func Sizeof(value reflect.Value) int {
	switch value.Kind() {
	case reflect.Slice:
		if s := sizeofByType(value.Type().Elem()); s >= 0 {
			return s * value.Len()
		}

		return -1

	case reflect.String:
		return value.Len() + 2 //string len + uint16

	case reflect.Struct:
		sum := 0
		size := 0
		for i, n := 0, value.NumField(); i < n; i++ {
			//fieldType := t.Field(i).Type.Kind()
			size = Sizeof(value.Field(i))

			// if fieldType == reflect.String {
			// 	size = Sizeof(value.Field(i)) //string len + uint16
			// } else {
			// 	size = sizeof(t.Field(i).Type)
			// }

			if size < 0 {
				return -1
			}

			sum += size
		}
		return sum

	default:
		return sizeofByType(value.Type())
	}
}

func sizeofByType(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Array:
		if s := sizeofByType(t.Elem()); s >= 0 {
			return s * t.Len()
		}
		return -1

	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return int(t.Size())
	}

	return -1
}
