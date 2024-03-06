package logger

type Field struct {
	key   string
	value interface{}
}

func Why(key string, value interface{}) Field {
	return Field{key: key, value: value}
}
