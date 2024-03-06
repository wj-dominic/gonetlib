package logger_test

import (
	"gonetlib/logger"
	"testing"
)

func TestLogger(t *testing.T) {
	mylogger := logger.Create(logger.Config{})

	mylogger.Debug("debug log", logger.Why("test", 10), logger.Why("test2", "this is test"))

}
