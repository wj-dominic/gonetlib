package logger

import (
	"strconv"
	"strings"
)

type Field struct {
	key   string
	value interface{}
}

func Why(key string, value interface{}) Field {
	return Field{key: key, value: value}
}

func (field *Field) ToString() string {
	var sb strings.Builder
	sb.WriteString(field.key)
	sb.WriteString(":")

	var value string

	switch field.value.(type) {
	case bool:
		_value := field.value.(bool)
		value = strconv.FormatBool(_value)
		break

	case int8:
		_value := field.value.(int8)
		value = strconv.Itoa(int(_value))
		break
	case int16:
		_value := field.value.(int16)
		value = strconv.Itoa(int(_value))
		break
	case int32:
		_value := field.value.(int32)
		value = strconv.Itoa(int(_value))
		break
	case int64:
		_value := field.value.(int64)
		value = strconv.FormatInt(_value, 10)
		break
	case int:
		_value := field.value.(int)
		value = strconv.Itoa(_value)
		break

	case uint8:
		_value := field.value.(uint8)
		value = strconv.FormatUint(uint64(_value), 10)
		break
	case uint16:
		_value := field.value.(uint16)
		value = strconv.FormatUint(uint64(_value), 10)
		break
	case uint32:
		_value := field.value.(uint32)
		value = strconv.FormatUint(uint64(_value), 10)
		break
	case uint64:
		_value := field.value.(uint64)
		value = strconv.FormatUint(_value, 10)
		break
	case uint:
		_value := field.value.(uint)
		value = strconv.FormatUint(uint64(_value), 10)
		break

	case float32:
		_value := field.value.(float32)
		value = strconv.FormatFloat(float64(_value), 'f', -1, 64)
		break

	case float64:
		_value := field.value.(float64)
		value = strconv.FormatFloat(_value, 'f', -1, 64)
		break

	case string:
		value = field.value.(string)
		break

	default:
		value = "invalid"
		break
	}

	sb.WriteString(value)

	return sb.String()
}
