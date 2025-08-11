package task

import (
	"fmt"

	"github.com/wj-dominic/gonetlib/util/snowflake"
)

type Task[Out any] interface {
	Start(...interface{}) Task[Out]
	Await(func(Out, error)) Task[Out]
	Wait() Task[Out]
	Result() (Out, error)
}

type gonetTask[Out any] struct {
	id          uint64
	job         func(...interface{}) (Out, error)
	await       func(Out, error)
	result      Out
	resultChan  chan Out
	error       error
	numOfThread uint8
}

func (t *gonetTask[Out]) Start(params ...interface{}) Task[Out] {
	f := func(threadId uint8) bool {
		defer func() {
			fmt.Printf("end async job | thread id %d\n", threadId)
		}()

		fmt.Printf("begin async job | thread id %d\n", threadId)

		result, err := t.job(params...)

		t.error = err
		t.resultChan <- result

		return true
	}

	add(f, t.numOfThread)
	return t
}

func (t *gonetTask[Out]) Await(await func(Out, error)) Task[Out] {
	t.await = await

	awaitFunc := func(threadId uint8) bool {
		defer func() {
			fmt.Printf("end await job | thread id %d\n", threadId)
		}()

		fmt.Printf("begin await job | thread id %d\n", threadId)

		select {
		case result, ok := <-t.resultChan:
			if ok == false {
				return false
			}

			t.await(result, t.error)
		default:
			return false
		}

		return true
	}

	add(awaitFunc)

	return t
}

func (t *gonetTask[Out]) Wait() Task[Out] {
	var ok bool
	t.result, ok = <-t.resultChan
	if ok == false {
		t.error = fmt.Errorf("cannot get result from channel")
	}

	return t
}

func (t *gonetTask[Out]) Result() (Out, error) {
	return t.result, t.error
}

func New[Out any](f func(...interface{}) (Out, error), numOfThread ...uint8) Task[Out] {
	tempNumOfThread := uint8(0)
	if len(numOfThread) != 0 {
		tempNumOfThread = numOfThread[0]
	}

	return &gonetTask[Out]{
		id:          snowflake.GenerateID(int64(tempNumOfThread)),
		job:         f,
		await:       nil,
		resultChan:  make(chan Out, 1),
		error:       nil,
		numOfThread: tempNumOfThread,
	}
}
