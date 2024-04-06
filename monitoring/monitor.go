package monitoring

import (
	"sync"
	"time"
)

type Collector[T any] interface {
	Collect() (T, error)
}

type Monitor struct {
	ticker    *time.Ticker
	collector Collector[interface{}]
	recvCount uint64
	interval  int

	prevSessionMonitoringData    SessionMonitoringData
	currentSessionMonitoringData SessionMonitoringData

	dataLock sync.RWMutex

	done chan struct{}
}

func NewMonitor() *Monitor {
	return &Monitor{
		recvCount: 1,
		interval:  1,
		done:      make(chan struct{}),
		dataLock:  sync.RWMutex{},
	}
}

func (m *Monitor) Start() error {
	m.ticker = time.NewTicker(time.Duration(m.interval) * time.Second)

	go func() {
		for {
			select {
			case <-m.ticker.C:
				m.dataLock.Lock()
				m.prevSessionMonitoringData = m.currentSessionMonitoringData

				monitoringData, _ := m.collector.Collect()
				m.currentSessionMonitoringData = monitoringData.(SessionMonitoringData)

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

func (m *Monitor) GetData() MonitoringDataResponse {
	m.dataLock.RLock()
	defer m.dataLock.RUnlock()

	resp := MonitoringDataResponse{
		ActiveSessions:      m.currentSessionMonitoringData.ActiveSessions,
		ConnectableSessions: m.currentSessionMonitoringData.ConnectableSessions,

		SendTPS: m.currentSessionMonitoringData.SendCount - m.prevSessionMonitoringData.SendCount,
		RecvTPS: m.currentSessionMonitoringData.RecvCount - m.prevSessionMonitoringData.RecvCount,

		SendBPS: m.currentSessionMonitoringData.SendBytes - m.prevSessionMonitoringData.SendBytes,
		RecvBPS: m.currentSessionMonitoringData.RecvBytes - m.prevSessionMonitoringData.RecvBytes,

		SendChannelCount: m.currentSessionMonitoringData.SendChannelCount,
	}

	return resp
}

func (m *Monitor) AddCollector(collector Collector[interface{}]) {
	m.collector = collector
}
