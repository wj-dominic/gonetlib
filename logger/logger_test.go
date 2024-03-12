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

	_logger.Dispose()
}

func TestLoggerWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	config := logger.CreateLoggerConfig()
	config.MinimumLevel(logger.DebugLevel)
	config.WriteToConsole()
	config.WriteToFile(logger.WriteToFile{
		Filepath:        "other_log.txt",
		RollingInterval: logger.RollingIntervalDay,
	})

	_logger := config.CreateLoggerWithContext(ctx)

	for i := 0; i < 100; i++ {
		_logger.Debug("debug log", logger.Why("count", i))
	}

	time.Sleep(time.Second)

	for i := 0; i < 100; i++ {
		_logger.Info("hello log", logger.Why("count", i), logger.Why("test", "is test"))
	}

	cancel()

	//wait를 위해 명시적으로 호출, 암시적으로 기다리게 할 수 있는 방법은 없을까?
	_logger.Dispose()
}
