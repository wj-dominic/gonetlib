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
	_logger.Debug("debug log", logger.Why("test", 10), logger.Why("test2", "this is test"))

	cancel()
}
