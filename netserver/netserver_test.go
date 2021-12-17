package netserver

import (
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	go Run("127.0.0.1:8888")
	time.Sleep(time.Second * 30)
}
