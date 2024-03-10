package logger_test

import (
	"context"
	"gonetlib/logger"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	config := logger.CreateLoggerConfig()
	config.MinimumLevel(logger.DebugLevel)
	config.WriteToConsole()
	config.WriteToFile(logger.WriteToFile{
		Filepath:        "log.txt",
		RollingInterval: logger.RollingIntervalDay,
	})

	_logger := config.CreateLogger(ctx)

	for i := 0; i < 100; i++ {
		_logger.Debug("debug log", logger.Why("count", i))
	}

	time.Sleep(time.Second)

	for i := 0; i < 100; i++ {
		_logger.Info("hello log", logger.Why("count", i), logger.Why("test", "is test"))
	}

	//TODO:종료 부분 손보기, context로 종료하는 것과 별도로 종료하는 것 2가지 타입으로 나누기
	cancel()
	_logger.Close()
}
