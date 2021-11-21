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
	Info
	Debug
	Max
)

var levelStr = [4]string{"ERROR", "WARNING", "INFO", "DEBUG"}

type Logger struct {
	logFile   *os.File
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
		logName = filepath.Base(os.Args[0])
		logName = logName[:len(logName)-len(filepath.Ext(logName))] + ".log"
	}

	msg := make(chan string)
	stop := make(chan bool)

	return &Logger{
		logFile:   nil,
		level:     level,
		directory: dir,
		logName:   logName,
		msg:       msg,
		stop:      stop,
	}
}

func (l *Logger) Start() error {
	err := l.setDirectory()
	if err != nil {
		return err
	}

	logFile, err := os.OpenFile(l.directory+"/"+l.logName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	l.logFile = logFile

	go l.loggerProc()

	return nil
}

func (l *Logger) Stop() {
	l.stop <- true
	l.logFile.Close()
}

func (l *Logger) Error(msg string) {
	l.log(Error, msg)
}

func (l *Logger) Warn(msg string) {
	l.log(Warning, msg)
}

func (l *Logger) Info(msg string) {
	l.log(Info, msg)
}

func (l *Logger) Debug(msg string) {
	l.log(Debug, msg)
}

func (l *Logger) log(level Level, msg string) {
	if level > l.level {
		return
	}

	log := fmt.Sprintf("[%s][%s]%s\n", time.Now().Format("2006-01-02 15:04:05"), levelStr[level], msg)
	l.msg <- log
}

func (l *Logger) loggerProc() {
	// TODO: ticker duration 변경 가능하도록
	ticker := time.NewTicker(time.Second * 3)

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
			l.logFile.WriteString(msg)
		default:
			return
		}
	}
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) setDirectory() error {
	if isExistFile(l.directory) {
		return nil
	}

	err := os.MkdirAll(l.directory, os.ModePerm)

	return err
}

func isExistFile(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}

	return true
}
