package logger

type ILogger interface {
	Debug(message string, fields ...Field)
	Info(message string, fields ...Field)
	Warn(message string, fields ...Field)
	Error(message string, fields ...Field)
}

type Logger struct {
}

func Create(config Config) ILogger {
	return &Logger{}
}

func (logger *Logger) Debug(message string, fields ...Field) {
}

func (logger *Logger) Info(message string, fields ...Field) {

}
func (logger *Logger) Warn(message string, fields ...Field) {

}
func (logger *Logger) Error(message string, fields ...Field) {
}
