package task

import (
	"fmt"
	"gonetlib/util/snowflake"
)

type ITask[Out any] interface {
	Start(...interface{}) ITask[Out]
	Await(func(Out, error)) ITask[Out]
	Wait() ITask[Out]
	Result() (Out, error)
}

type Task[Out any] struct {
	id          uint64
	job         func(...interface{}) (Out, error)
	await       func(Out, error)
	result      Out
	resultChan  chan Out
	error       error
	numOfThread uint8
}

func (t *Task[Out]) Start(params ...interface{}) ITask[Out] {
	f := func(threadId uint8) bool {
		defer func() {
			fmt.Printf("end job | thread id %d\n", threadId)
		}()

		fmt.Printf("begin job | thread id %d\n", threadId)

		result, err := t.job(params...)

		t.error = err
		t.resultChan <- result

		return true
	}

	add(f, t.numOfThread)
	return t
}

func (t *Task[Out]) Await(await func(Out, error)) ITask[Out] {
	t.await = await

	awaitFunc := func(threadId uint8) bool {
		defer func() {
			//fmt.Printf("end await job | thread id %d\n", threadId)
		}()

		//fmt.Printf("begin await job | thread id %d\n", threadId)

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

func (t *Task[Out]) Wait() ITask[Out] {
	var ok bool
	t.result, ok = <-t.resultChan
	if ok == false {
		t.error = fmt.Errorf("cannot get result from channel")
	}

	return t
}

func (t *Task[Out]) Result() (Out, error) {
	return t.result, t.error
}

func New[Out any](f func(...interface{}) (Out, error), numOfThread ...uint8) ITask[Out] {
	tempNumOfThread := uint8(0)
	if len(numOfThread) != 0 {
		tempNumOfThread = numOfThread[0]
	}

	return &Task[Out]{
		id:          snowflake.GenerateID(int64(tempNumOfThread)),
		job:         f,
		await:       nil,
		resultChan:  make(chan Out, 1),
		error:       nil,
		numOfThread: tempNumOfThread,
	}
}
