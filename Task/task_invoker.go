package task

import (
	"context"
	"gonetlib/netlogger"
)

func NewTaskInvoker(invokerID uint8, bucket *chan ITask) *TaskInvoker {
	return &TaskInvoker{
		id:     invokerID,
		bucket: bucket,
	}
}

type TaskInvoker struct {
	id     uint8
	bucket *chan ITask

	stopFunc context.CancelFunc
}

func (i *TaskInvoker) Run() bool {
	ctx, cancelFunc := context.WithCancel(context.Background())
	i.stopFunc = cancelFunc
	go i.proc(ctx)
	return true
}

func (i *TaskInvoker) Stop() bool {
	i.stopFunc()
	return true
}

func (i *TaskInvoker) proc(ctx context.Context) {
	for {
		select {
		case task := <-*i.bucket:
			task.Run()
			break

		case _ = <-ctx.Done():
			netlogger.GetLogger().Debug("task invoker is done | id[%d]", i.id)
			return
		}
	}
}
