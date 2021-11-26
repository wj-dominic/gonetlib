package util

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	lock	sync.RWMutex
	done 	uint32
}

func (once *Once) Do(fn func()) {
	if atomic.LoadUint32(&once.done) == 1 {
		return
	}

	once.lock.Lock()
	defer once.lock.Unlock()
	if once.done == 0 {
		once.done = 1
		fn()
	}
}

func (once *Once) Reset() {
	if atomic.LoadUint32(&once.done) == 0{
		return
	}

	once.lock.Lock()
	defer once.lock.Unlock()
	once.done = 0
}
