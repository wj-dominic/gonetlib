package stopwatch_test

import (
	"fmt"
	"gonetlib/util/stopwatch"
	"testing"
	"time"
)

func TestStopwatch(t *testing.T) {
	watch := stopwatch.New()
	watch.Start()
	time.Sleep(time.Second * 5)
	watch.Stop()

	fmt.Println(watch.ElapsedTime())
}
