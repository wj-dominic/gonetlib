package monitoring

import (
	"net/http"
)

type Exporter interface {
	Start() error
	Stop()
}

type defaultExporter struct {
	server  *http.Server
	monitor *Monitor
}

func NewDefaultExporter(monitor *Monitor) Exporter {
	return &defaultExporter{
		monitor: monitor,
	}
}

func (de *defaultExporter) Start() error {
	de.server = &http.Server{Addr: ":8080"}
	http.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := de.monitor.MonitoringData()
		w.Write(resp)
	})

	go func() {
		de.server.ListenAndServe()
	}()

	return nil
}

func (de *defaultExporter) Stop() {
	de.server.Close()
}