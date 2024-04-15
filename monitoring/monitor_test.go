package monitoring

import (
	"gonetlib/logger"
	"math/rand"
	"testing"
	"time"
)

type SampleMonitoringData struct {
	SendCount uint64
	RecvCount uint64

	SendBytes uint64
	RecvBytes uint64
}

type TestCollector struct {
	sampleData SampleMonitoringData
}

func (c *TestCollector) Collect() (interface{}, error) {
	c.sampleData.SendCount += rand.Uint64() % 10
	c.sampleData.RecvCount += rand.Uint64() % 10

	c.sampleData.SendBytes += rand.Uint64() % 10
	c.sampleData.RecvBytes += rand.Uint64() % 10

	return c.sampleData, nil
}

func TestStart(t *testing.T) {
	config := logger.CreateLoggerConfig().
		WriteToConsole().
		WriteToFile(
			logger.WriteToFile{
				Filepath:        "./test_monitoring.log",
				RollingInterval: logger.RollingIntervalDay,
			}).
		MinimumLevel(logger.DebugLevel).
		TickDuration(1000)
	_logger := config.CreateLogger()

	m := NewMonitor(_logger)
	e := NewDefaultExporter(m)
	testCollector := &TestCollector{}
	m.AddCollector(testCollector)

	err := m.Start()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	e.Start()
	// Wait for some time to allow ticks to occur
	time.Sleep(20 * time.Second)

	e.Stop()
}
