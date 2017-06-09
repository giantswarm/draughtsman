package slack

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
	prometheusSubsystem = "slack_notifier"
)

var (
	requestDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "request_duration_milliseconds",
			Help:      "Time taken to make Slack requests.",
		},
	)
	requestTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "request_total",
			Help:      "Number of Slack requests.",
		},
	)
)

func init() {
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestTotal)
}

func updateSlackMetrics(startTime time.Time) {
	requestDuration.Set(float64(time.Since(startTime) / time.Millisecond))
	requestTotal.Inc()
}
