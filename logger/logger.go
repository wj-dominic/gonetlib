package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type ILogger interface {
	Debug(message string, fields ...Field)
	Info(message string, fields ...Field)
	Warn(message string, fields ...Field)
	Error(message string, fields ...Field)
	Close()
}

var wg sync.WaitGroup

type Logger struct {
	config          config
	logs            chan string
	writeToFileLogs chan string
	ctx             context.Context
}

func CreateLogger(config config, ctx context.Context) ILogger {
	logger := &Logger{
		config:          config,
		logs:            make(chan string),
		writeToFileLogs: make(chan string),
		ctx:             ctx,
	}

	wg.Add(2)
	go logger.tick()
	go logger.writeToFile()

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

func (logger *Logger) Close() {
	wg.Wait()

	close(logger.logs)
	close(logger.writeToFileLogs)
}

func (logger *Logger) log(level Level, message string, fields ...Field) {
	//레벨이 낮으면 로그 안씀
	if logger.config.limitLevel > level {
		return
	}

	//문자열 조합
	var sb strings.Builder
	sb.WriteString(message)

	for _, field := range fields {
		sb.WriteString("|")
		sb.WriteString(field.ToString())
	}

	//채널에 삽입
	logger.logs <- sb.String()
}

func (logger *Logger) tick() {
	for {
		select {
		case log := <-logger.logs:
			logger.write(log)
		case <-logger.ctx.Done():
			fmt.Println("tick is end")
			wg.Done()
			return
		}
	}
}

func (logger *Logger) write(log string) {
	if logger.config.writeToConsole.enable == true {
		fmt.Println(log)
	}

	if logger.config.writeToFile.enable == true {
		logger.writeToFileLogs <- log
	}
}

func (logger *Logger) writeToFile() {
	var sb strings.Builder
	ontick := time.NewTicker(logger.config.tickDuration)
	for {
		select {
		case <-ontick.C:
			err := flushToFile(&logger.config.writeToFile, sb.String())
			if err != nil {
				panic(err)
			}

			sb.Reset()

		case log := <-logger.writeToFileLogs:
			//로그 모으기
			sb.WriteString(log)
			sb.WriteString("\n")

		case <-logger.ctx.Done():
			if len(logger.writeToFileLogs) > 0 {
				for log := range logger.writeToFileLogs {
					sb.WriteString(log)
					sb.WriteString("\n")
				}
			}

			err := flushToFile(&logger.config.writeToFile, sb.String())
			if err != nil {
				panic(err)
			}

			sb.Reset()

			fmt.Println("writeToFile is end")
			wg.Done()
			return
		}
	}
}

func flushToFile(wtf *WriteToFile, text string) error {
	if wtf == nil {
		return fmt.Errorf("invalid object of writeToFile")
	}

	if len(text) == 0 {
		return nil
	}

	//rolling interval 기준으로 파일 이름 구하기
	path := wtf.makeRollingFilepath()

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
