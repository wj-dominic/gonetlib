package logger_test

import (
	"gonetlib/logger"
	"testing"
)

func TestLogger(t *testing.T) {

	config := logger.CreateLoggerConfig()
	config.SetLimitLevel(logger.DebugLevel)
	config.WriteToConsole()
	config.WriteToFile(logger.WriteToFile{
		Filepath:        "log.txt",
		RollingInterval: logger.RollingIntervalDay,
	})

	_logger := config.CreateLogger()
	_logger.Debug("debug log", logger.Why("test", 10), logger.Why("test2", "this is test"))
}
