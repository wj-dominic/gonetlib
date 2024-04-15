package session

type MonitoringData struct {
	ActiveSessions      uint64
	ConnectableSessions uint64

	BySession map[uint64]SessionMonitoringData
}

type SessionMonitoringData struct {
	SendCount uint64
	RecvCount uint64

	SendBytes uint64
	RecvBytes uint64

	SendChannelCount uint64
}
