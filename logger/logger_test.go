package logger

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logger := NewLogger(Max, "", "")
	t.Logf("\ndir: %s\nfile name: %s\n", logger.directory, logger.logName)

	err := logger.Start()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		id := i
		wg.Add(1)
		go func() {
			for i := 0; i < 10; i++ {
				if i > 0 && i%5 == 0 {
					time.Sleep(time.Millisecond * 100)
				}
				msg := fmt.Sprintf("[id:%d]test log %d", id, i)
				logger.Debug(msg)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
