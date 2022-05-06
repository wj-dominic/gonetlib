package task

import "gonetlib/util/singleton"

type TaskBucket struct {
	ID uint16

	buckets  []chan ITask
	invokers []IInvoker
}

func NewBucket(bucketID uint16, maxInvokerSize, maxTaskSize uint8) *TaskBucket {
	bucket := new(TaskBucket)

	//TODO :: Bucket 만들기

}

func Bucket(bucketID uint16) *TaskBucket {
	bucketManager := singleton.GetInstance[TaskBucketManager]()
	return bucketManager.buckets[bucketID]
}
