package monitoring

type SessionMonitoringData struct {
	ActiveSessions      uint64
	ConnectableSessions uint64

	SendCount uint64
	RecvCount uint64

	SendBytes uint64
	RecvBytes uint64

	SendChannelCount uint64
}

type MonitoringDataResponse struct {
	ActiveSessions      uint64 `json:"activeSessions"`
	ConnectableSessions uint64 `json:"connectableSessions"`

	SendTPS uint64 `json:"sendTPS"`
	RecvTPS uint64 `json:"recvTPS"`

	SendBPS uint64 `json:"sendBPS"`
	RecvBPS uint64 `json:"recvBPS"`

	SendChannelCount uint64 `json:"sendChannelCount"`
}

func (smd *SessionMonitoringData) Add(other SessionMonitoringData) {
	smd.SendCount += other.SendCount
	smd.RecvCount += other.RecvCount

	smd.SendBytes += other.SendBytes
	smd.RecvBytes += other.RecvBytes
}
