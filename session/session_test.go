package session_test

import (
	"fmt"
	"gonetlib/util"
	"sync"
	"testing"
	"unsafe"
)

type Data struct {
	RefCount    int32
	ReleaseFlag int32
}

func TestSession(t *testing.T) {
	data := Data{RefCount: 0, ReleaseFlag: 0}
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			util.InterlockIncrement(&data.RefCount)
		}()
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			refCount := util.InterlockDecrement(&data.RefCount)
			if refCount == 0 {
				origin := (*int64)(unsafe.Pointer(&data))
				exchange := (*int64)(unsafe.Pointer(&Data{0, 1}))
				compare := (*int64)(unsafe.Pointer(&Data{0, 0}))

				if util.InterlockedCompareExchange64(origin, *exchange, *compare) == true {
					fmt.Println("going to release!!")
				}
			}
		}()
	}

	wg.Wait()
}
