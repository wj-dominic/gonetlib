package logger_test

import (
	"context"
	"gonetlib/logger"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	config := logger.CreateLoggerConfig()
	config.MinimumLevel(logger.InfoLevel)
	config.WriteToConsole()
	config.WriteToFile(logger.WriteToFile{
		Filepath:        "log.txt",
		RollingInterval: logger.RollingIntervalDay,
	})

	_logger := config.CreateLogger()

	for i := 0; i < 100; i++ {
		_logger.Debug("debug log", logger.Why("count", i))
	}

	time.Sleep(time.Second)

	for i := 0; i < 100; i++ {
		_logger.Info("hello log", logger.Why("count", i), logger.Why("test", "is test"))
	}

	_logger.Close()
}

func TestLoggerWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	config := logger.CreateLoggerConfig()
	config.MinimumLevel(logger.DebugLevel)
	config.WriteToConsole()
	config.WriteToFile(logger.WriteToFile{
		Filepath:        "other_log.txt",
		RollingInterval: logger.RollingIntervalDay,
		RollingFileSize: 100,
	})

	_logger := config.CreateLoggerWithContext(ctx)

	for i := 0; i < 100; i++ {
		_logger.Debug("debug log", logger.Why("count", i))
	}

	time.Sleep(time.Second)

	for i := 0; i < 100; i++ {
		_logger.Info("hello log", logger.Why("count", i), logger.Why("test", "is test"))
	}

	//TODO:외부에서 캔슬하면 로거는 어떻게 Wait???
	cancel()
}
