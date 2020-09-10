package snorlax

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(latencyHist)
}

// latencyHist measures each request's latency.
var latencyHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "snorlax",
	Subsystem: "requests",
	Name:      "latency",
	Help:      "Request latency in seconds",
}, []string{"method", "code", "path"})
