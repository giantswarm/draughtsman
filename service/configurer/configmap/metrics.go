package configmap

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// prometheusNamespace is the namespace to use for Prometheus metrics.
	// See: https://godoc.org/github.com/prometheus/client_golang/prometheus#Opts
	prometheusNamespace = "draughtsman"

	// prometheusSubsystem is the subsystem to use for Prometheus metrics.
	// See: https://godoc.org/github.com/prometheus/client_golang/prometheus#Opts
	prometheusSubsystem = "configmap_configurer"
)

var (
	requestDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "request_duration_milliseconds",
			Help:      "Time taken to request ConfigMap manifests from Kubernetes.",
		},
	)
	requestTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "request_total",
			Help:      "Number of requests to fetch ConfigMap manifests from Kubernetes.",
		},
	)
)

func init() {
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestTotal)
}

func updateConfigMapMetrics(startTime time.Time) {
	requestDuration.Set(float64(time.Since(startTime) / time.Millisecond))
	requestTotal.Inc()
}
