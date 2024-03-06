package logger

import "fmt"

type Level uint8

const (
	DebugLevel Level = 0 + iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

func (level Level) ToString() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return fmt.Sprintf("other(%d)", level)
	}
}
