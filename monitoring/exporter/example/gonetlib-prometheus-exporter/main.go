package main

import (
	"gonetlib-prometheus-exporter/exporter"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

func main() {
	prometheus.Register(version.NewCollector("gonetlib-exporter"))
	prometheus.Register(exporter.NewExporter())

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h := promhttp.HandlerFor(prometheus.Gatherers{
			prometheus.DefaultGatherer,
		}, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
