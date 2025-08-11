package logger_test

import (
	"testing"
	"time"

	"github.com/wj-dominic/gonetlib/logger"
)

func TestLogger(t *testing.T) {
	config := logger.NewLoggerConfig()
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
	config := logger.NewLoggerConfig()
	config.MinimumLevel(logger.DebugLevel)
	config.WriteToConsole()
	config.WriteToFile(logger.WriteToFile{
		Filepath:        "other_log.txt",
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

	time.Sleep(time.Second)

	for i := 0; i < 100; i++ {
		_logger.Info("hello 2 log")
	}

	_logger.Dispose()
}

func TestDefaultLogger(t *testing.T) {
	for i := 0; i < 100; i++ {
		logger.Debug("debug log", logger.Why("count", i))
	}

	time.Sleep(time.Second)

	for i := 0; i < 100; i++ {
		logger.Info("hello log", logger.Why("count", i), logger.Why("test", "is test"))
	}

	logger.Dispose()
}
