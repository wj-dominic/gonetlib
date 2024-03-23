package snowflake

import (
	"os"
	"sync"
	"time"
)

const (
	epoch int64 = 1288834974657
)

var (
	lock     sync.Mutex
	lastTime int64 = 0
	sequence int64 = 0
)

func GenerateID(machineID int64) uint64 {
	// defer lock.Unlock()
	// lock.Lock()

	generatedID := uint64(0)

	now := time.Now().UnixMilli()

	if lastTime == now {
		sequence = (sequence + 1) & 4095
		if sequence == 0 {
			for now <= lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		sequence = 0
	}

	lastTime = now

	pid := (int64(os.Getpid()) % 31)

	generatedID = (uint64)(((now - epoch) << 22) | (machineID << 17) | (pid << 12) | (sequence))
	return generatedID
}
