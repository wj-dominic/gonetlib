package task

import (
	"context"
	"gonetlib/netlogger"
	"sync"
)

func NewTaskInvoker(invokerID uint8, bucket *chan ITask) *TaskInvoker {
	return &TaskInvoker{
		id:     invokerID,
		bucket: bucket,
	}
}

type TaskInvoker struct {
	id       uint8
	bucket   *chan ITask
	wg       sync.WaitGroup
	stopFunc context.CancelFunc
}

func (i *TaskInvoker) Run() bool {
	ctx, cancelFunc := context.WithCancel(context.Background())
	i.stopFunc = cancelFunc
	i.wg.Add(1)
	go i.proc(ctx)
	return true
}

func (i *TaskInvoker) Stop() bool {
	i.stopFunc()
	i.wg.Wait()
	return true
}

func (i *TaskInvoker) proc(ctx context.Context) {
	for {
		select {
		case task := <-*i.bucket:
			task.Run()

		case <-ctx.Done():
			netlogger.Debug("task invoker is done | id[%d]", i.id)
			return
		}
	}
}
