package monitoring

import (
	"encoding/json"
	"gonetlib/logger"
	"sync"
	"time"
)

type Collector interface {
	Collect() (interface{}, error)
}

type Monitor struct {
	ticker    *time.Ticker
	collector Collector
	interval  int

	monitoringData []byte

	dataLock sync.RWMutex

	done chan struct{}

	logger logger.ILogger
}

func NewMonitor(logger logger.ILogger) *Monitor {
	return &Monitor{
		interval: 1,
		done:     make(chan struct{}),
		dataLock: sync.RWMutex{},
		logger:   logger,
	}
}

func (m *Monitor) Start() error {
	m.ticker = time.NewTicker(time.Duration(m.interval) * time.Second)

	go func() {
		for {
			select {
			case <-m.ticker.C:
				if m.collector == nil {
					m.logger.Error("Collector is not set")
					continue
				}

				m.dataLock.Lock()

				coll, err := m.collector.Collect()
				if err != nil {
					m.logger.Error("Error collecting monitoring data:", logger.Why("error", err))
					m.dataLock.Unlock()
					continue
				}

				m.monitoringData, err = json.Marshal(coll)
				if err != nil {
					m.logger.Error("Error marshalling monitoring data:", logger.Why("error", err))
					m.dataLock.Unlock()
					continue
				}

				m.dataLock.Unlock()
			case <-m.done:
				m.ticker.Stop()
				return
			}
		}
	}()

	return nil
}

func (m *Monitor) Stop() {
	m.done <- struct{}{}
}

func (m *Monitor) MonitoringData() []byte {
	m.dataLock.RLock()
	defer m.dataLock.RUnlock()

	return m.monitoringData
}

func (m *Monitor) AddCollector(collector Collector) {
	m.collector = collector
}
