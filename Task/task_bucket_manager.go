package task

import "gonetlib/netlogger"

type TaskBucketManager struct {
	buckets map[uint16]*TaskBucket
}

func (m *TaskBucketManager) Init() {
	m.buckets = make(map[uint16]*TaskBucket)
}

func (m *TaskBucketManager) CreateBucket(bucketID uint16, maxInvokerSize uint8, maxTaskSize uint8) bool {
	if _, exist := m.buckets[bucketID]; exist == true {
		netlogger.GetLogger().Warn("already has bucket | id[%d]", bucketID)
		return false
	}

	bucket := NewBucket(bucketID, maxInvokerSize, maxTaskSize)
	if bucket == nil {
		netlogger.GetLogger().Warn("Failed to create new bucket | id[%d]", bucketID)
		return false
	}

	gym := NewGym(gymName, gymType)

	if gym.Create(trainerCount, routinesCount) == false {
		GetLogger().Error("cannot create a gym")
		gym = nil
		return false
	}

	gymManager.gyms[gymType] = gym

	return true
}
