package metrics

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var ScanPortsTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "scan_ports_total",
	Help: "Total number of ports scanned",
})

var ScanPortsOpenTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "scan_ports_open_total",
	Help: "Total number of open ports found",
})

var ScanErrors = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "scan_errors_total",
	Help: "Total number of scan errors",
})

var ExportErrors = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "export_errors_total",
	Help: "Total number of JSON export errors",
})

var ScanDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "scan_duration_seconds",
	Help:    "Duration of port scans",
	Buckets: prometheus.DefBuckets,
})

func init() {
	prometheus.MustRegister(ScanPortsOpenTotal, ExportErrors, ScanErrors, ScanDuration, ScanPortsTotal)
}

func StartServer(addr string) {
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			slog.Error("metrics server failed", "error", err)
		}
	}()

}
