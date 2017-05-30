package eventer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// githubDeploymentApiUrlTemplate is the string template for the GitHub
	// API call to fetch Deployments.
	githubDeploymentApiUrlTemplate = "https://api.github.com/repos/%s/%s/deployments"

	// etagHeader is the header used for etag.
	// See: https://en.wikipedia.org/wiki/HTTP_ETag.
	etagHeader = "Etag"

	// rateLimitLimitHeader is the header set by GitHub to show the total
	// rate limit value.
	// See: https://developer.github.com/v3/#rate-limiting
	rateLimitLimitHeader = "X-RateLimit-Limit"

	// rateLimitRemainingHeader is the header set by GitHub to show the
	// remaining rate limit value.
	// See: https://developer.github.com/v3/#rate-limiting
	rateLimitRemainingHeader = "X-RateLimit-Remaining"

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
)

func init() {
	prometheus.MustRegister(rateLimitLimit)
	prometheus.MustRegister(rateLimitRemaining)
	prometheus.MustRegister(deploymentRequestDuration)
	prometheus.MustRegister(deploymentResponseCodeTotal)
}

// githubDeployment represents a GitHub API Deployment.
// See: https://developer.github.com/v3/repos/deployments/#get-a-single-deployment
type githubDeployment struct {
	// Sha is the SHA hash of the commit the deployment references.
	Sha string `json:"sha"`

	// Environment is the environment field of the Github deployment.
	Environment string `json:"environment"`
}

// DeploymentEvent returns the githubDeployment as a DeploymentEvent.
func (g githubDeployment) DeploymentEvent(project string) DeploymentEvent {
	return DeploymentEvent{
		Name: project,
	}
}

// GithubEventer is an Eventer that uses Github Deployment Events as a backend.
var GithubEventer EventerType = "GithubEventer"

// githubEventer is an implementer of the Eventer interface,
// that uses GitHub Deployment Events as a backend.
type githubEventer struct {
	// Dependencies.
	client *http.Client
	logger micrologger.Logger

	// Settings.
	environment  string
	oauthToken   string
	organisation string
	pollInterval time.Duration
	projectList  []string
}

// fetchNewDeploymentEvents fetches any new GitHub Deployment Events for the
// given project.
// It uses Etags and If-None-Match headers to avoid requesting the same events
// multiple times, to avoid hitting rate limits.
// See: https://developer.github.com/v3/#conditional-requests
func (e *githubEventer) fetchNewDeploymentEvents(project string, etagMap map[string]string) ([]githubDeployment, error) {
	e.logger.Log("debug", "fetching deployment events", "project", project)

	startTime := time.Now()

	deploymentUrl := fmt.Sprintf(
		githubDeploymentApiUrlTemplate,
		e.organisation,
		project,
	)

	req, err := http.NewRequest("GET", deploymentUrl, nil)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", e.oauthToken))

	// If we have an etag header for this project, then we have already
	// requested deployment events for it.
	// So, set the header, so we only get notified of new events.
	if val, ok := etagMap[project]; ok {
		req.Header.Set("If-None-Match", val)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}
	defer resp.Body.Close()

	// Update metrics.
	rateLimitLimitValue, err := strconv.ParseFloat(resp.Header.Get(rateLimitLimitHeader), 64)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}
	rateLimitLimit.Set(rateLimitLimitValue)

	rateLimitRemainingValue, err := strconv.ParseFloat(resp.Header.Get(rateLimitRemainingHeader), 64)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}
	rateLimitRemaining.Set(rateLimitRemainingValue)

	deploymentRequestDuration.WithLabelValues(e.organisation, project, strconv.Itoa(resp.StatusCode)).Set((float64(time.Since(startTime) / time.Millisecond)))
	deploymentResponseCodeTotal.WithLabelValues(e.organisation, project, strconv.Itoa(resp.StatusCode)).Inc()

	// If there are no new deployment events, return quickly.
	if resp.StatusCode == http.StatusNotModified {
		e.logger.Log("debug", "no new deployment events, continuing", "project", project)
		return []githubDeployment{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		e.logger.Log("error", "received a non-200 status code", "code", resp.StatusCode)
		return nil, microerror.MaskAnyf(unexpectedStatusCode, fmt.Sprintf("received non-200 status code: %v", resp.StatusCode))
	}

	// Save the new etag header, so we don't get these deployment events again.
	etagMap[project] = resp.Header.Get(etagHeader)

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	var deployments []githubDeployment
	if err := json.Unmarshal(bytes, &deployments); err != nil {
		return nil, microerror.MaskAny(err)
	}

	if len(deployments) > 0 {
		e.logger.Log("debug", "found new deployment events", "project", project)
	}

	// Filter out deployment events that do not match the current environment.
	matches := []githubDeployment{}

	for _, deployment := range deployments {
		if deployment.Environment == e.environment {
			matches = append(matches, deployment)
		}
	}

	return matches, nil
}

func (e *githubEventer) NewDeploymentEvents() (<-chan DeploymentEvent, error) {
	e.logger.Log("debug", "starting polling for github deployment events", "interval", e.pollInterval)

	deploymentEventChannel := make(chan DeploymentEvent)
	ticker := time.NewTicker(e.pollInterval)

	go func() {
		etagMap := make(map[string]string)

		for c := ticker.C; ; <-c {
			for _, project := range e.projectList {
				deployments, err := e.fetchNewDeploymentEvents(project, etagMap)
				if err != nil {
					e.logger.Log("error", "could not fetch deployment events", "message", err.Error())
				}

				for _, deployment := range deployments {
					deploymentEventChannel <- deployment.DeploymentEvent(project)
				}
			}
		}
	}()

	return deploymentEventChannel, nil
}

func (e *githubEventer) SetPending(event DeploymentEvent) error {
	return nil
}

func (e *githubEventer) SetSuccess(event DeploymentEvent) error {
	return nil
}

func (e *githubEventer) SetFailed(event DeploymentEvent) error {
	return nil
}
