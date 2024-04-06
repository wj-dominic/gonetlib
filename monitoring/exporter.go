package monitoring

import (
	"encoding/json"
	"gonetlib/logger"
	"net/http"
)

type Exporter interface {
	Run() error
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

func (de *defaultExporter) Run() error {
	de.server = &http.Server{Addr: ":8080"}
	http.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data, _ := json.Marshal(de.monitor.GetData())
		logger.Info("monitoring data", logger.Why("data", string(data)))
		w.Write(data)
	})

	go func() {
		de.server.ListenAndServe()
	}()

	return nil
}

func (de *defaultExporter) Stop() {
	de.server.Close()
}
