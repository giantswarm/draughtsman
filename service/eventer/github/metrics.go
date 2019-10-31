package github

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/giantswarm/microerror"
)

const (
	// prometheusNamespace is the namespace to use for Prometheus metrics.
	// See: https://godoc.org/github.com/prometheus/client_golang/prometheus#Opts
	prometheusNamespace = "draughtsman"

	// prometheusSubsystem is the subsystem to use for Prometheus metrics.
	// See: https://godoc.org/github.com/prometheus/client_golang/prometheus#Opts
	prometheusSubsystem = "github_eventer"
)

var (
	rateLimitLimit = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: prometheusNamespace,
		Subsystem: prometheusSubsystem,
		Name:      "rate_limit_limit",
		Help:      "Rate limit limit for GitHub API requests.",
	})
	rateLimitRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: prometheusNamespace,
		Subsystem: prometheusSubsystem,
		Name:      "rate_limit_remaining",
		Help:      "Rate limit remaining for GitHub API requests.",
	})

	deploymentRequestDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "github_deployment_duration_milliseconds",
			Help:      "Time taken to request GitHub deployments.",
		},
		[]string{"organisation", "project", "code"},
	)
	deploymentResponseCodeTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "github_deployment_response_code",
			Help:      "Response codes of GitHub API requests for deployments.",
		},
		[]string{"organisation", "project", "code"},
	)

	deploymentStatusRequestDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "github_deployment_status_duration_milliseconds",
			Help:      "Time taken to request GitHub deployment statuses.",
		},
		[]string{"method", "organisation", "project", "code"},
	)
	deploymentStatusResponseCodeTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "github_deployment_status_response_code",
			Help:      "Response codes of GitHub API requests for deployment statuses.",
		},
		[]string{"method", "organisation", "project", "code"},
	)
)

func init() {
	prometheus.MustRegister(rateLimitLimit)
	prometheus.MustRegister(rateLimitRemaining)

	prometheus.MustRegister(deploymentRequestDuration)
	prometheus.MustRegister(deploymentResponseCodeTotal)

	prometheus.MustRegister(deploymentStatusRequestDuration)
	prometheus.MustRegister(deploymentStatusResponseCodeTotal)
}

// updateRateLimitMetrics is a utility function that takes a Response
// containing rate limit headers, and updates the rate limit metrics.
func updateRateLimitMetrics(response *http.Response) error {
	rateLimitLimitValue, err := parseRateLimitValue(response)
	if err != nil {
		return microerror.Mask(err)
	}
	rateLimitLimit.Set(rateLimitLimitValue)

	rateLimitRemainingValue, err := parseRateLimitRemaining(response)
	if err != nil {
		return microerror.Mask(err)
	}
	rateLimitRemaining.Set(rateLimitRemainingValue)

	return nil
}

// updateDeploymentMetrics is a utility function for updating metrics related
// to Deployment API calls.
func updateDeploymentMetrics(organisation, project string, statusCode int, startTime time.Time) {
	deploymentRequestDuration.WithLabelValues(
		organisation,
		project,
		strconv.Itoa(statusCode),
	).Set(
		float64(time.Since(startTime) / time.Millisecond),
	)

	deploymentResponseCodeTotal.WithLabelValues(
		organisation,
		project,
		strconv.Itoa(statusCode),
	).Inc()
}

// updateDeploymentStatusMetrics is a utility function for updating metrics
// related to Deployment Status API calls.
func updateDeploymentStatusMetrics(method, organisation, project string, statusCode int, startTime time.Time) {
	deploymentStatusRequestDuration.WithLabelValues(
		method,
		organisation,
		project,
		strconv.Itoa(statusCode),
	).Set(
		float64(time.Since(startTime) / time.Millisecond),
	)

	deploymentStatusResponseCodeTotal.WithLabelValues(
		method,
		organisation,
		project,
		strconv.Itoa(statusCode),
	).Inc()
}
