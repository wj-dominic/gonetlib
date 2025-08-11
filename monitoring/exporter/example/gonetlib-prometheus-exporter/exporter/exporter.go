package exporter

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "gonetlib"
)

var sample = `
{
    "active_sessions": 0,
    "connectable_sessions": 0,
    "BySession": {
        "1784952844399570944": {
            "send_count": 0,
            "recv_count": 0,
            "send_bytes": 0,
            "recv_bytes": 0,
            "send_channel_count": 0
        },
        "1784952844454096896": {
            "send_count": 0,
            "recv_count": 0,
            "send_bytes": 0,
            "recv_bytes": 0,
            "send_channel_count": 0
        },
		"17849528444540962396": {
            "send_count": 0,
            "recv_count": 0,
            "send_bytes": 0,
            "recv_bytes": 0,
            "send_channel_count": 0
        },
		"2784952844454096896": {
            "send_count": 0,
            "recv_count": 0,
            "send_bytes": 0,
            "recv_bytes": 0,
            "send_channel_count": 0
        },
		"3784952844454096896": {
            "send_count": 0,
            "recv_count": 0,
            "send_bytes": 0,
            "recv_bytes": 0,
            "send_channel_count": 0
        }
    }
}`

type Exporter struct {
	dataReqClient *http.Client

	activeSessions      *prometheus.Desc
	connectableSessions *prometheus.Desc

	sendCount        *prometheus.Desc
	recvCount        *prometheus.Desc
	sendBytes        *prometheus.Desc
	recvBytes        *prometheus.Desc
	sendChannelCount *prometheus.Desc
}

func NewExporter() *Exporter {
	return &Exporter{
		dataReqClient: &http.Client{},

		activeSessions: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "active_sessions"),
			"Current number of active sessions",
			[]string{},
			nil,
		),

		connectableSessions: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "connectable_sessions"),
			"Current number of connectable sessions",
			[]string{},
			nil,
		),

		sendCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "send_count"),
			"Number of sent messages",
			[]string{"sid"},
			nil,
		),

		recvCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "recv_count"),
			"Number of received messages",
			[]string{"sid"},
			nil,
		),

		sendBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "send_bytes"),
			"Number of sent bytes",
			[]string{"sid"},
			nil,
		),

		recvBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "recv_bytes"),
			"Number of received bytes",
			[]string{"sid"},
			nil,
		),

		sendChannelCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "send_channel_count"),
			"Number of messages in the send channel",
			[]string{"sid"},
			nil,
		),
	}
}

func (e *Exporter) Describe(descs chan<- *prometheus.Desc) {
	descs <- e.activeSessions
	descs <- e.connectableSessions
	descs <- e.sendCount
	descs <- e.recvCount
	descs <- e.sendBytes
	descs <- e.recvBytes
	descs <- e.sendChannelCount
}

func (e *Exporter) Collect(metrics chan<- prometheus.Metric) {
	// sample to data
	data := make(map[string]interface{})
	err := json.NewDecoder(strings.NewReader(sample)).Decode(&data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data)

	// active_sessions
	metrics <- prometheus.MustNewConstMetric(e.activeSessions, prometheus.GaugeValue, rand.Float64()*100)

	// connectable_sessions
	metrics <- prometheus.MustNewConstMetric(e.connectableSessions, prometheus.GaugeValue, rand.Float64()*100)

	// BySession
	bySession, _ := data["BySession"].(map[string]interface{})
	for sid, _ := range bySession {
		metrics <- prometheus.MustNewConstMetric(e.sendCount, prometheus.GaugeValue, rand.Float64()*100, sid)
		metrics <- prometheus.MustNewConstMetric(e.recvCount, prometheus.GaugeValue, rand.Float64()*100, sid)
		metrics <- prometheus.MustNewConstMetric(e.sendBytes, prometheus.GaugeValue, rand.Float64()*100, sid)
		metrics <- prometheus.MustNewConstMetric(e.recvBytes, prometheus.GaugeValue, rand.Float64()*100, sid)
		metrics <- prometheus.MustNewConstMetric(e.sendChannelCount, prometheus.GaugeValue, rand.Float64()*100, sid)
	}
}
