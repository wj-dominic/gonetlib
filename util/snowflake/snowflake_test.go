package snowflake_test

import (
	"fmt"
	"gonetlib/util/snowflake"
	"sync"
	"testing"
)

func TestGenerateID(t *testing.T) {
	var wg sync.WaitGroup
	var container sync.Map

	//중복이 있다!?
	for i := 0; i < 1500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			genID := snowflake.GenerateID(1)
			_, loaded := container.LoadOrStore(genID, "test")
			if loaded == true {
				fmt.Println(genID)
			}
		}()
	}

	wg.Wait()

	count := int64(0)
	container.Range(func(k, v interface{}) bool {
		count++
		return true
	})

	fmt.Println(count)
}
