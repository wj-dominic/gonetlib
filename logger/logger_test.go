package logger

import (
	"fmt"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logger := NewLogger(Debug, "", "")
	t.Logf("\ndir: %s\nfile name: %s\n", logger.directory, logger.logName)

	logger.Start()
	for i := 0; i < 50; i++ {
		if i > 0 && i%10 == 0 {
			time.Sleep(time.Second)
		}
		msg := fmt.Sprintf("test%d", i)
		logger.Log(msg)
	}
	logger.Stop()
}
