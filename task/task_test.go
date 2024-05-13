package task_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/wj-dominic/gonetlib/logger"
	"github.com/wj-dominic/gonetlib/task"
)

func TestTaskSummation(t *testing.T) {

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()

			sumation := task.New(
				func(params ...interface{}) (int, error) {
					time.Sleep(time.Second * 5)
					return params[0].(int) + params[1].(int) + params[2].(int), nil
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

	task.Dispose()
	logger.Dispose()
}

func TaskMain(i ...interface{}) (error, error) {
	something := task.New(func(i ...interface{}) (int, error) {
		//do something long task...
		time.Sleep(time.Second * 5)
		return i[0].(int) + i[1].(int), nil
	}, 1)

	something.Start(1, 3).Await(func(result int, err error) {
		if err != nil {
			fmt.Printf("something error, error %s\n", err.Error())
			return
		}

		fmt.Printf("task return is %d\n", result)
	})

	return nil, nil
}

func TestTaskCallback(t *testing.T) {
	taskMain := task.New(TaskMain)
	taskMain.Start()
	taskMain.Wait()
	time.Sleep(time.Second * 15)
}
