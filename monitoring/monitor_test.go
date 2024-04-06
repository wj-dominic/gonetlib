package monitoring

import (
	"math/rand"
	"testing"
	"time"
)

type TestCollector struct {
	sampleData SessionMonitoringData
}

func (c *TestCollector) Collect() (interface{}, error) {
	c.sampleData.SendCount += rand.Uint64() % 10
	c.sampleData.RecvCount += rand.Uint64() % 10

	c.sampleData.SendBytes += rand.Uint64() % 10
	c.sampleData.RecvBytes += rand.Uint64() % 10

	return c.sampleData, nil
}

func TestStart(t *testing.T) {
	m := NewMonitor()
	e := NewDefaultExporter(m)
	testCollector := &TestCollector{}
	m.AddCollector(testCollector)

	err := m.Start()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	e.Run()

	// Wait for some time to allow ticks to occur
	time.Sleep(20 * time.Second)

	m.Stop()
	e.Stop()
}
