package netserver

import (
	"gonetlib/logger"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	s := NewNetServer()
	s.Start()
	time.Sleep(time.Millisecond * 500)

	l := logger.GetLogger()
	l.Start()
	l.Error("TestLogger")

	s.Stop()
}
