package logger

import (
	"strings"
	"time"
)

const (
	DELIM string = "|"
)

type Log struct {
	time    string
	level   Level
	message string
	fields  []Field
}

func NewLog(level Level, message string, fields ...Field) Log {
	now := time.Now()
	return Log{
		time:    now.Format("2006-01-02 15:04:05.000"),
		level:   level,
		message: message,
		fields:  fields,
	}
}

func (log *Log) ToString() string {
	var sb strings.Builder
	sb.WriteString(log.time)
	sb.WriteString("|")
	sb.WriteString(log.level.ToString())
	sb.WriteString("|")
	sb.WriteString(log.message)

	for _, field := range log.fields {
		sb.WriteString("|")
		sb.WriteString(field.ToString())
	}

	return sb.String()
}
