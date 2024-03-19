package task

import (
	"gonetlib/netlogger"
	"gonetlib/util/singleton"
)

func GetBucket(bucketID uint16) *TaskBucket {
	bucketManager := singleton.GetInstance[TaskBucketManager]()
	if _, exist := bucketManager.buckets[bucketID]; exist == false {
		if bucketManager.CreateBucket(bucketID, 1, 1) == false {
			return nil
		}
	}

	return bucketManager.buckets[bucketID]
}

type TaskBucketManager struct {
	buckets map[uint16]*TaskBucket
}

func (m *TaskBucketManager) Init() {
	m.buckets = make(map[uint16]*TaskBucket)
}

func (m *TaskBucketManager) CreateBucket(bucketID uint16, maxInvokers uint8, maxBuckets uint8) bool {
	if _, exist := m.buckets[bucketID]; exist == true {
		netlogger.Warning("already has bucket | id[%d]", bucketID)
		return false
	}

	bucket := NewBucket(bucketID, maxInvokers, maxBuckets)
	if bucket == nil {
		netlogger.Warning("Failed to create new bucket | id[%d]", bucketID)
		return false
	}

	m.buckets[bucketID] = bucket

	return true
}
