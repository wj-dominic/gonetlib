package logger

import (
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	DELIM string = "||"
)

type Log struct {
	time         string
	functionName string
	line         string
	level        Level
	message      string
	fields       []Field
}

func NewLog(level Level, message string, fields ...Field) Log {
	now := time.Now()
	var functionName, lineStr string

	if level >= ErrorLevel {
		_, fn, lineInt, _ := runtime.Caller(3)
		functionName = fn
		lineStr = strconv.Itoa(lineInt)
	}

	return Log{
		time:         now.Format("2006-01-02 15:04:05.000"),
		functionName: functionName,
		line:         lineStr,
		level:        level,
		message:      message,
		fields:       fields,
	}
}

func (log *Log) ToString() string {
	var sb strings.Builder
	sb.WriteString(log.time)
	sb.WriteString(DELIM)

	sb.WriteString(log.level.ToString())
	sb.WriteString(DELIM)
	sb.WriteString(log.message)

	for _, field := range log.fields {
		sb.WriteString(DELIM)
		sb.WriteString(field.ToString())
	}

	if log.level >= ErrorLevel {
		sb.WriteString(DELIM)
		sb.WriteString(log.functionName)
		sb.WriteString(":")
		sb.WriteString(log.line)
	}

	return sb.String()
}
