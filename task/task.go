package task

import (
	"fmt"
	"gonetlib/util/snowflake"
)

type ITask[Out any] interface {
	Start(...interface{}) ITask[Out]
	Wait() ITask[Out]
	Result() (Out, error)
}

type Task[Out any] struct {
	id          uint64
	job         func(...interface{}) Out
	result      Out
	resultChan  chan Out
	error       error
	numOfThread uint8
}

func (t *Task[Out]) Start(params ...interface{}) ITask[Out] {
	f := func(threadId uint8) {
		fmt.Printf("begin job | thread id %d\n", threadId)
		t.resultChan <- t.job(params...)
		fmt.Printf("end job | thread id %d\n", threadId)
	}

	add(f, t.numOfThread)
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

func New[Out any](f func(...interface{}) Out, numOfThread ...uint8) ITask[Out] {
	tempNumOfThread := uint8(0)
	if len(numOfThread) != 0 {
		tempNumOfThread = numOfThread[0]
	}

	return &Task[Out]{
		id:          snowflake.GenerateID(int64(tempNumOfThread)),
		job:         f,
		resultChan:  make(chan Out, 1),
		error:       nil,
		numOfThread: tempNumOfThread,
	}
}
