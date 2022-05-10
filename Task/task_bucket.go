package task

import "gonetlib/netlogger"

const (
	MaxTasks = 300
)

func NewBucket(bucketID uint16, maxInvokers, maxBuckets uint8) *TaskBucket {
	bucket := new(TaskBucket)

	bucket.ID = bucketID

	bucket.buckets = make([]chan ITask, maxBuckets)
	bucket.invokers = make([]IInvoker, maxInvokers)

	for index := range bucket.buckets {
		bucket.buckets[index] = make(chan ITask, MaxTasks)
	}

	for index := range bucket.invokers {
		invokerID := uint8(index)
		bucketID := invokerID % maxBuckets

		bucket.invokers[index] = NewTaskInvoker(invokerID, &bucket.buckets[bucketID])

		bucket.invokers[index].Run()
	}

	return bucket
}

type TaskBucket struct {
	ID uint16

	buckets  []chan ITask
	invokers []IInvoker
}

func (b *TaskBucket) AddTask(task ITask, invokerID uint16) bool {
	if task == nil {
		netlogger.GetLogger().Error("Invalid task")
		return false
	}

	if b.invokers[invokerID] == nil {
		netlogger.GetLogger().Error("Not found invoker | invokerID[%d] invokers[%d]", invokerID, len(b.invokers))
		return false
	}

	b.buckets[invokerID] <- task

	return true
}
