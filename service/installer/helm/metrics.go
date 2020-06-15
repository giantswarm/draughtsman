package helm

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	gaugeValue float64 = 1

	// prometheusNamespace is the namespace to use for Prometheus metrics.
	// See: https://godoc.org/github.com/prometheus/client_golang/prometheus#Opts
	prometheusNamespace = "draughtsman"

	// prometheusSubsystem is the subsystem to use for Prometheus metrics.
	// See: https://godoc.org/github.com/prometheus/client_golang/prometheus#Opts
	prometheusSubsystem = "helm_installer"
)

var (
	helmCommandDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "helm_command_duration_milliseconds",
			Help:      "Time taken to execute Helm commands.",
		},
		[]string{"name"},
	)
	helmCommandTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "helm_command_total",
			Help:      "Number of total Helm commands run.",
		},
		[]string{"name"},
	)
	helmReleaseFailure = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "helm_release_failed",
			Help:      "the number of failed releases for draughtsman deployments",
		},
		[]string{"name", "status"},
	)
)

func init() {
	prometheus.MustRegister(helmCommandDuration)
	prometheus.MustRegister(helmCommandTotal)
	prometheus.MustRegister(helmReleaseFailure)
}

// updateHelmMetrics is a utility function for updating metrics related to
// Helm command calls.
func updateHelmMetrics(name string, startTime time.Time) {
	helmCommandDuration.WithLabelValues(
		name,
	).Set(
		float64(time.Since(startTime) / time.Millisecond),
	)

	helmCommandTotal.WithLabelValues(
		name,
	).Inc()
}

func reportHelmRelease(name, status string) {
	helmReleaseFailure.WithLabelValues(
		name, status,
	).Set(gaugeValue)
}
