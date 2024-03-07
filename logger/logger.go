package logger

import "strings"

type ILogger interface {
	Debug(message string, fields ...Field)
	Info(message string, fields ...Field)
	Warn(message string, fields ...Field)
	Error(message string, fields ...Field)
}

type Logger struct {
	config config

	logs chan string
}

func CreateLogger(config config) ILogger {
	return &Logger{
		config: config,
		logs:   make(chan string),
	}
}

func (logger *Logger) Debug(message string, fields ...Field) {
	logger.write(DebugLevel, message, fields...)
}
func (logger *Logger) Info(message string, fields ...Field) {
	logger.write(InfoLevel, message, fields...)
}
func (logger *Logger) Warn(message string, fields ...Field) {
	logger.write(WarnLevel, message, fields...)
}
func (logger *Logger) Error(message string, fields ...Field) {
	logger.write(ErrorLevel, message, fields...)
}

func (logger *Logger) write(level Level, message string, fields ...Field) {
	if logger.config.limitLevel > level {
		return
	}

	//문자열 조합
	var sb strings.Builder
	sb.WriteString(message)

	for _, field := range fields {
		sb.WriteString("|")
		sb.WriteString(field.ToString())
	}

	//채널에 삽입
	logger.logs <- sb.String()
}
