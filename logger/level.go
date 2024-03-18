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
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return fmt.Sprintf("OTHER(%d)", level)
	}
}
