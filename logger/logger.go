package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Level uint32

const (
	Error = 0 + iota
	Warning
	Debug
	Notice
	Max
)

type Logger struct {
	level     Level
	directory string
	logName   string
	msg       chan string
	stop      chan bool
}

func NewLogger(level Level, dir string, logName string) *Logger {
	if len(dir) == 0 {
		dir = "./"
	}

	if len(logName) == 0 {
		fileName := os.Args[0]
		logName = fileName[:len(fileName)-len(filepath.Ext(fileName))] + ".log"
	}

	msg := make(chan string)

	stop := make(chan bool)

	return &Logger{
		level:     Error,
		directory: dir,
		logName:   logName,
		msg:       msg,
		stop:      stop,
	}
}

func (l *Logger) Start() {
	go l.loggerProc()
}

func (l *Logger) Stop() {
	l.stop <- true
}

func (l *Logger) Log(msg string) {
	l.msg <- msg
}

func (l *Logger) loggerProc() {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-l.stop:
			return
		case <-ticker.C:
			l.writeLog()
		}
	}
}

func (l *Logger) writeLog() {
	for {
		select {
		case msg := <-l.msg:
			fmt.Println(msg)
		default:
			return
		}
	}
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) SetDirectory(dir string) {
	l.directory = dir
}

func (l *Logger) SetLogName(logName string) {
	l.logName = logName
}
