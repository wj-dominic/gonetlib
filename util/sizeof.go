package util

import (
	"reflect"
)

func Sizeof(valueType reflect.Type) int {
	switch valueType.Kind() {
	case reflect.Array, reflect.Slice:
		if elemSize := Sizeof(valueType.Elem()) ; elemSize >= 0 {
			return elemSize * valueType.Len()
		}
		break

	case reflect.Struct:
		sum := 0
		for idx, max := 0, valueType.NumField() ; idx < max ; idx++ {
			fieldSize := Sizeof(valueType.Field(idx).Type)
			if fieldSize < 0 {
				return -1
			}
			sum += fieldSize
		}
		return sum

	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return int(valueType.Size())

	case reflect.Ptr, reflect.Uintptr:
		return int(valueType.Size())
	}

	return -1
}