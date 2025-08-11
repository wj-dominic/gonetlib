package task

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wj-dominic/gonetlib/util"
)

type taskRunner struct {
	Ids        sync.Map
	jobs       []chan func(uint8) bool
	wg         sync.WaitGroup
	isDisposed int32
	inUse      int32
}

var runner *taskRunner = newTaskRunner(uint8(runtime.NumCPU()))

func newTaskRunner(maxCount uint8) *taskRunner {
	runner := &taskRunner{
		Ids:        sync.Map{},
		jobs:       make([]chan func(uint8) bool, maxCount),
		wg:         sync.WaitGroup{},
		isDisposed: 0,
	}

	for i := range runner.jobs {
		runner.jobs[i] = make(chan func(uint8) bool, 100)

		runner.wg.Add(1)
		go func(id uint8) {
			defer runner.wg.Done()

			goroutineId := getGoroutineId()
			if goroutineId == -1 {
				panic("failed to get gorouine id")
			}

			runner.Ids.Store(goroutineId, id)

			ontick := time.NewTicker(time.Second * 3)

		Loop:
			for {
				select {
				case job, ok := <-runner.jobs[id]:
					if ok == false {
						break Loop
					}
					if job(id) == false {
						runner.jobs[id] <- job
					}
				case <-ontick.C:
					if runner.isDisposed == 1 {
						break Loop
					}
				}
				time.Sleep(time.Millisecond)
			}

			//flush
			if len(runner.jobs[id]) > 0 {
				for job := range runner.jobs[id] {
					job(id)
				}
			}

		}(uint8(i))
	}

	return runner
}

func getGoroutineId() int {
	buf := make([]byte, 32)
	n := runtime.Stack(buf, false)
	buf = buf[:n]

	text := string(buf)
	splits := strings.Split(text, " ")
	goId, err := strconv.Atoi(splits[1])
	if err != nil {
		return -1
	}

	return goId
}

func add(f func(uint8) bool, numOfThread ...uint8) error {
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

	tempNumOfThread := getRunnerId(getGoroutineId())
	if len(numOfThread) != 0 {
		tempNumOfThread = numOfThread[0]
	}

	if len(runner.jobs) <= int(tempNumOfThread) {
		tempNumOfThread = tempNumOfThread % uint8(len(runner.jobs))
	}

	runner.jobs[tempNumOfThread] <- f

	return nil
}

func getRunnerId(goroutineId int) uint8 {
	runnerId, ok := runner.Ids.Load(goroutineId)
	if ok == false {
		return 0
	}

	return runnerId.(uint8)
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
