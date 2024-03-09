package netlogger

import (
	"fmt"
	"gonetlib/util"
	"gonetlib/util/singleton"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

type Level uint32

const (
	ErrorLevel Level = 0 + iota
	WarningLevel
	InfoLevel
	DebugLevel
	MaxLevel
)

var levelStr = [4]string{"ERROR", "WARNING", "INFO", "DEBUG"}

type defaultNetLogger struct {
	loggers   [2]*NetLogger
	curIndex  uint32
	prevIndex uint32

	logOption *NetLoggerOption
}

func (dl *defaultNetLogger) current() *NetLogger {
	return dl.loggers[dl.curIndex]
}

func (dl *defaultNetLogger) prev() *NetLogger {
	return dl.loggers[dl.prevIndex]
}

func (dl *defaultNetLogger) change() {
	dl.prev().SetOption(dl.logOption)
	dl.prevIndex = atomic.SwapUint32(&dl.curIndex, dl.prevIndex)
}

func (dl *defaultNetLogger) setoption(option *NetLoggerOption) {
	dl.logOption = option
	dl.change()
}

var std = &defaultNetLogger{loggers: [2]*NetLogger{New(), New()}, curIndex: 0, prevIndex: 1, logOption: newLoggerOption()}

func logging(level Level, format string, args ...interface{}) {
	std.current().Log(level, format, args...)
}

func Error(format string, args ...interface{}) {
	logging(ErrorLevel, format, args...)
}

func Warning(format string, args ...interface{}) {
	logging(WarningLevel, format, args...)
}

func Info(format string, args ...interface{}) {
	logging(InfoLevel, format, args...)
}

func Debug(format string, args ...interface{}) {
	logging(WarningLevel, format, args...)
}

func SetOption(option *NetLoggerOption) *NetLogger {
	std.setoption(option)
	return std.current()
}

func SetLevel(level Level) *NetLogger {
	std.logOption.SetLevel(level)
	std.setoption(std.logOption)
	return std.current()
}

func SetFileName(name string) *NetLogger {
	std.logOption.SetFileName(name)
	std.setoption(std.logOption)
	return std.current()
}

func SetTickDuration(millisec time.Duration) *NetLogger {
	std.logOption.SetTickDuration(millisec)
	std.setoption(std.logOption)
	return std.current()
}

type NetLogger struct {
	logFile   *os.File
	msg       chan string
	stop      chan bool
	isRunning uint32

	option *NetLoggerOption

	wg sync.WaitGroup
}

type NetLoggerOption struct {
	level        Level
	tickDuration time.Duration
	logFileName  string
}

func newLoggerOption() *NetLoggerOption {
	return &NetLoggerOption{
		level:        MaxLevel,
		tickDuration: time.Second * 3,
		logFileName:  "./netlogger",
	}
}

func (op *NetLoggerOption) SetLevel(level Level) {
	op.level = level
}

func (op *NetLoggerOption) SetFileName(name string) {
	if len(name) > 0 {
		dir := filepath.Dir(name)
		if util.IsExistFile(dir) == false {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return
			}
		}

		if len(filepath.Ext(name)) >= 0 {
			baseName := filepath.Base(name)
			name = dir + "/" + baseName
		}

		op.logFileName = name
	}
}

func (op *NetLoggerOption) SetTickDuration(millisec time.Duration) {
	op.tickDuration = millisec
}

func New() *NetLogger {
	var logger NetLogger
	logger.Init()

	return &logger
}

func GetLogger() *NetLogger {
	return singleton.GetInstance[NetLogger]()
}

func (l *NetLogger) Init() {
	l.logFile = nil
	l.msg = make(chan string)
	l.stop = make(chan bool)

	l.option = newLoggerOption()
}

func (l *NetLogger) Start() error {
	if atomic.CompareAndSwapUint32(&l.isRunning, 0, 1) == false {
		return fmt.Errorf("already start a logger")
	}

	nowTime := time.Now().Format("20060102150405")
	fileName := l.option.logFileName + "_" + nowTime + ".log"
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	l.logFile = logFile

	l.wg.Add(1)

	go l.Tick()

	return nil
}

func (l *NetLogger) Stop() bool {
	if atomic.CompareAndSwapUint32(&l.isRunning, 1, 0) == false {
		return false
	}

	l.stop <- true

	l.wg.Wait()

	l.logFile.Close()

	return true
}

func (l *NetLogger) Error(format string, v ...interface{}) {
	l.Log(ErrorLevel, format, v...)
}

func (l *NetLogger) Warn(format string, v ...interface{}) {
	l.Log(WarningLevel, format, v...)
}

func (l *NetLogger) Info(format string, v ...interface{}) {
	l.Log(InfoLevel, format, v...)
}

func (l *NetLogger) Debug(format string, v ...interface{}) {
	l.Log(DebugLevel, format, v...)
}

func (l *NetLogger) Log(level Level, format string, v ...interface{}) {
	if l.isRunning == 0 {
		if err := l.Start(); err != nil {
			return
		}
	}

	if level > l.option.level {
		return
	}

	msg := fmt.Sprintf(format, v...)
	log := fmt.Sprintf("[%s][%s]%s \n", time.Now().Format("2006-01-02 15:04:05"), levelStr[level], msg)
	l.msg <- log
}

func (l *NetLogger) Tick() {
	defer l.wg.Done()

	ticker := time.NewTicker(l.option.tickDuration)

	for {
		select {
		case <-l.stop:
			l.flushLog()
			return
		case <-ticker.C:
			l.flushLog()
		}
	}
}

func (l *NetLogger) flushLog() {
	for {
		select {
		case msg := <-l.msg:
			l.writeLog(msg)
		default:
			return
		}
	}
}

func (l *NetLogger) writeLog(msg string) {
	l.logFile.WriteString(msg)
}

func (l *NetLogger) SetOption(option *NetLoggerOption) *NetLogger {
	l.SetLevel(option.level).SetFileName(option.logFileName).SetTickDuration(option.tickDuration)
	return l
}

func (l *NetLogger) SetLevel(level Level) *NetLogger {
	l.option.SetLevel(level)
	return l
}

func (l *NetLogger) SetFileName(name string) *NetLogger {
	l.option.SetFileName(name)
	return l
}

func (l *NetLogger) SetTickDuration(millisec time.Duration) *NetLogger {
	l.option.SetTickDuration(millisec)
	return l
}
