package logger

import (
	"fmt"
	"gonetlib/util"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type ILogger interface {
	Debug(message string, fields ...Field)
	Info(message string, fields ...Field)
	Warn(message string, fields ...Field)
	Error(message string, fields ...Field)
	Dispose()
}

type Logger struct {
	config     config
	logs       chan Log
	wg         sync.WaitGroup
	isDisposed int32
}

func CreateLogger(config config) ILogger {
	logger := &Logger{
		config:     config,
		logs:       make(chan Log),
		isDisposed: 0,
	}

	logger.wg.Add(1)
	go logger.tick()

	return logger
}

func (logger *Logger) Debug(message string, fields ...Field) {
	logger.log(DebugLevel, message, fields...)
}
func (logger *Logger) Info(message string, fields ...Field) {
	logger.log(InfoLevel, message, fields...)
}
func (logger *Logger) Warn(message string, fields ...Field) {
	logger.log(WarnLevel, message, fields...)
}
func (logger *Logger) Error(message string, fields ...Field) {
	logger.log(ErrorLevel, message, fields...)
}

func (logger *Logger) Dispose() {
	if util.InterlockedCompareExchange(&logger.isDisposed, 1, 0) == false {
		return
	}

	close(logger.logs)
	logger.wg.Wait()
}

func (logger *Logger) log(level Level, message string, fields ...Field) {
	if atomic.LoadInt32(&logger.isDisposed) == 1 {
		return
	}

	//레벨이 낮으면 로그 안씀
	if logger.config.limitLevel > level {
		return
	}

	//채널에 삽입
	logger.logs <- NewLog(level, message, fields...)
}

func (logger *Logger) tick() {
	defer logger.wg.Done()

	var sb strings.Builder
	toWriteFile := time.NewTicker(logger.config.tickDuration)

out:
	for {
		select {
		case log, ok := <-logger.logs:
			if ok == false {
				break out
			}

			//콘솔에는 즉시 출력
			if logger.config.writeToConsole.enable == true {
				fmt.Println(log.ToString())
			}

			//파일에는 모아서 출력
			if logger.config.writeToFile.enable == true {
				sb.WriteString(log.ToString())
				sb.WriteString("\n")
			}

		case <-toWriteFile.C:
			err := logger.flushToFile(sb.String())
			if err != nil {
				panic(err)
			}

			sb.Reset()
		}
	}

	//종료 시점에 남은게 있으면 파일에 쓰기
	if len(sb.String()) > 0 {
		err := logger.flushToFile(sb.String())
		if err != nil {
			panic(err)
		}
		sb.Reset()
	}
}

func (logger *Logger) flushToFile(text string) error {
	if len(text) == 0 {
		return nil
	}

	//rolling interval 기준으로 파일 이름 구하기
	path := logger.config.writeToFile.makeRollingFilepath()

	//해당 파일이 없으면 생성하기, 있으면 Append
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		return err
	}

	//모아놓은 로그 한번에 쓰기
	file.WriteString(text)
	file.Close()

	return nil
}

var _logger ILogger = CreateLoggerConfig().
	MinimumLevel(DebugLevel).
	WriteToConsole().
	WriteToFile(WriteToFile{
		Filepath:        "log.txt",
		RollingInterval: RollingIntervalDay,
	}).
	CreateLogger()

func Debug(message string, fields ...Field) {
	_logger.Debug(message, fields...)
}

func Info(message string, fields ...Field) {
	_logger.Info(message, fields...)
}

func Warn(message string, fields ...Field) {
	_logger.Warn(message, fields...)
}

func Error(message string, fields ...Field) {
	_logger.Error(message, fields...)
}

func Dispose() {
	_logger.Dispose()
}
