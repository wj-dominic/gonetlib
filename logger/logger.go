package logger

import (
	"context"
	"fmt"
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
	ontick := time.NewTicker(logger.config.tickDuration)
	ctx, cancel := context.WithTimeout(logger.ctx, logger.config.tickDuration)
	for {
		select {
		case <-ontick.C:
			logger.flush(ctx)
		case <-logger.ctx.Done():
			cancel()
			wg.Done()
		}
	}
}

func (logger *Logger) flush(ctx context.Context) {
	for {
		select {
		case log := <-logger.logs:
			logger.write(log)
		case <-ctx.Done():
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
	for {
		select {
		case <-logger.ctx.Done():
			wg.Done()
			return

		case log := <-logger.writeToFileLogs:
			//rolling interval 기준으로 파일 이름 구하기

			//해당 파일이 없으면 생성하기, 있으면 읽기

			//TODO:파일 쓰기
		}
	}
}
