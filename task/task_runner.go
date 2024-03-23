package task

import (
	"fmt"
	"gonetlib/util"
	"runtime"
	"sync"
	"time"
)

type taskRunner struct {
	jobs       []chan func(uint8)
	wg         sync.WaitGroup
	isDisposed int32
	inUse      int32
}

var runner *taskRunner = newTaskRunner(uint8(runtime.NumCPU()))

func newTaskRunner(maxCount uint8) *taskRunner {
	runner := &taskRunner{
		jobs:       make([]chan func(uint8), maxCount),
		wg:         sync.WaitGroup{},
		isDisposed: 0,
	}

	for i := range runner.jobs {
		runner.jobs[i] = make(chan func(uint8))

		runner.wg.Add(1)
		go func(id uint8) {
			defer runner.wg.Done()

			ontick := time.NewTicker(time.Second * 3)

		Loop:
			for {
				select {
				case job, ok := <-runner.jobs[id]:
					if ok == false {
						break Loop
					}
					job(id)
				case <-ontick.C:
					if runner.isDisposed == 1 {
						break Loop
					}
				}
			}

		}(uint8(i))
	}

	return runner
}

func add(f func(uint8), numOfThread ...uint8) error {
	defer func() {
		inUse := util.InterlockDecrement(&runner.inUse)
		if inUse == 0 && runner.isDisposed == 1 {
			Dispose()
		}
	}()

	util.InterlockIncrement(&runner.inUse)

	if runner.isDisposed == 1 {
		return fmt.Errorf("runner was disposed")
	}

	tempNumOfThread := uint8(0)
	if len(numOfThread) != 0 {
		tempNumOfThread = numOfThread[0]
	}

	if len(runner.jobs) <= int(tempNumOfThread) {
		tempNumOfThread = tempNumOfThread % uint8(len(runner.jobs))
	}

	runner.jobs[tempNumOfThread] <- f

	return nil
}

func Dispose() {
	if util.InterlockedCompareExchange(&runner.isDisposed, 1, 0) == false {
		return
	}

	if runner.inUse > 0 {
		return
	}

	for i, jobCh := range runner.jobs {
		if len(jobCh) > 0 {
			for job := range jobCh {
				job(uint8(i))
			}
		}

		close(jobCh)
	}
}
