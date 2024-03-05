package logger

type ILogger interface {
	Debug(string)
	Information(string)
	Warning(string)
	Error(string)
}

type Logger struct {
}

func (logger *Logger) Debug(message string) {
}

func (logger *Logger) Information(message string) {

}
func (logger *Logger) Warning(message string) {

}
func (logger *Logger) Error(message string) {

}
