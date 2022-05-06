package netlogger

import (
	"fmt"
	"gonetlib/util/singleton"
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

type NetLogger struct {
	logFile   *os.File
	level     Level
	directory string
	logName   string
	msg       chan string
	stop      chan bool
	isRunning bool
}

func (l *NetLogger) Init() {
	msg := make(chan string)
	stop := make(chan bool)

	l.logFile = nil
	l.level = Error
	l.directory = "./"
	l.logName = "Logger.log"
	l.msg = msg
	l.stop = stop
}

func GetLogger() *NetLogger {
	return singleton.GetInstance[NetLogger]()
}

func (l *NetLogger) Start() error {
	if l.isRunning {
		return nil
	}

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

	l.isRunning = true

	return nil
}

func (l *NetLogger) Stop() {
	l.stop <- true
	l.logFile.Close()
}

func (l *NetLogger) SetLogConfig(level Level, dir string, logName string) {
	l.level = level

	if len(dir) == 0 {
		dir = "./"
	}
	l.directory = dir

	if len(logName) == 0 {
		logName = filepath.Base(os.Args[0])
		logName = logName[:len(logName)-len(filepath.Ext(logName))] + ".log"
	}
	l.logName = logName
}

func (l *NetLogger) Error(format string, v ...interface{}) {
	l.log(Error, format, v...)
}

func (l *NetLogger) Warn(format string, v ...interface{}) {
	l.log(Warning, format, v...)
}

func (l *NetLogger) Info(format string, v ...interface{}) {
	l.log(Info, format, v...)
}

func (l *NetLogger) Debug(format string, v ...interface{}) {
	l.log(Debug, format, v...)
}

//func (l *NetLogger) log(level Level, msg string) {
func (l *NetLogger) log(level Level, format string, v ...interface{}) {
	if !l.isRunning {
		return
	}

	if level > l.level {
		return
	}

	msg := fmt.Sprintf(format, v...)
	log := fmt.Sprintf("[%s][%s]%s \n", time.Now().Format("2006-01-02 15:04:05"), levelStr[level], msg)
	l.msg <- log
}

func (l *NetLogger) loggerProc() {
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

func (l *NetLogger) writeLog() {
	for {
		select {
		case msg := <-l.msg:
			l.logFile.WriteString(msg)
		default:
			return
		}
	}
}

func (l *NetLogger) SetLevel(level Level) {
	l.level = level
}

func (l *NetLogger) setDirectory() error {
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
