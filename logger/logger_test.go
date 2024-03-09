package logger_test

import (
	"context"
	"gonetlib/logger"
	"testing"
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

	for i := 0; i < 10; i++ {
		_logger.Debug("debug log", logger.Why("count", i))
	}

	cancel()
	_logger.Close()
}
