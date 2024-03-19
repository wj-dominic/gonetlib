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
	var functionName, lineStr string = "", ""

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

	// ErrorLevel 이상일 때만 함수명과 라인번호를 출력
	// Info를 가장 낮추고 Debug부터 출력하는 것은 어떨지?
	if log.level >= ErrorLevel {
		sb.WriteString(log.functionName)
		sb.WriteString(":")
		sb.WriteString(log.line)
		sb.WriteString(DELIM)

	}

	sb.WriteString(log.level.ToString())
	sb.WriteString(DELIM)
	sb.WriteString(log.message)

	for _, field := range log.fields {
		sb.WriteString(DELIM)
		sb.WriteString(field.ToString())
	}

	return sb.String()
}
