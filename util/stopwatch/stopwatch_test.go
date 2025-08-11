package stopwatch_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/wj-dominic/gonetlib/util/stopwatch"
)

func TestStopwatch(t *testing.T) {
	watch := stopwatch.New()
	watch.Start()
	time.Sleep(time.Second * 5)
	watch.Stop()

	fmt.Println(watch.ElapsedTime())
}
