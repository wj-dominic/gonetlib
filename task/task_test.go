package task_test

import (
	"gonetlib/logger"
	"gonetlib/task"
	"sync"
	"testing"
	"time"
)

func TestTaskSummation(t *testing.T) {

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()

			sumation := task.New(
				func(params ...interface{}) int {
					time.Sleep(time.Second * 5)
					return params[0].(int) + params[1].(int) + params[2].(int)
				}, uint8(num))

			logger.Debug("begin start task", logger.Why("id", num))
			sumation.Start(num, 2, 3)
			sumation.Wait()
			logger.Debug("end start task", logger.Why("id", num))

			result, _ := sumation.Result()

			logger.Debug("task is done", logger.Why("id", num), logger.Why("result", result))
		}(i)
	}

	wg.Wait()

	logger.Dispose()
}
