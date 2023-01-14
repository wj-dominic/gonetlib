package netlogger

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		id := i
		wg.Add(1)
		go func() {
			for i := 0; i < 10; i++ {
				if i > 0 && i%2 == 0 {
					SetOption(&NetLoggerOption{MaxLevel, time.Second * 3, fmt.Sprintf("test_%d_%d", id, i)})
					time.Sleep(time.Millisecond * 1)
				}
				// msg := fmt.Sprintf("[id:%d]test log %d", id, i)
				Error("id:%d]test log %d", id, i)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestLogger(t *testing.T) {
	logger := GetLogger()
	logger.SetLevel(MaxLevel)
	t.Logf("file name: %s\n", logger.option.logFileName)

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
				// msg := fmt.Sprintf("[id:%d]test log %d", id, i)
				logger.Error("[id:%d]test log %d", id, i)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	logger.Stop()
}
