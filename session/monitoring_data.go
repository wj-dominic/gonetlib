package session

type MonitoringData struct {
	ActiveSessions      uint64 `json:"active_sessions"`
	ConnectableSessions uint64 `json:"connectable_sessions"`

	BySession map[uint64]SessionMonitoringData
}

type SessionMonitoringData struct {
	SendCount uint64 `json:"send_count"`
	RecvCount uint64 `json:"recv_count"`

	SendBytes uint64 `json:"send_bytes"`
	RecvBytes uint64 `json:"recv_bytes"`

	SendChannelCount uint64 `json:"send_channel_count"`
}
