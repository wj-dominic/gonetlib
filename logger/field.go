package logger

type Field struct {
	key   string
	value interface{}
}

func Why(key string, value interface{}) Field {
	return Field{key: key, value: value}
}

func (field *Field) ToString() string {
	//TODO:키, 밸류 조합해서 스트링 반환하기
	return ""
}
