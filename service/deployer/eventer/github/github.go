package github

import (
	"net/http"
	"time"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/draughtsman/service/deployer/eventer/spec"
)

// GithubEventer is an Eventer that uses Github Deployment Events as a backend.
var GithubEventer spec.EventerType = "GithubEventer"

// Config represents the configuration used to create a GitHub Eventer.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	Environment       string
	HTTPClientTimeout time.Duration
	OAuthToken        string
	Organisation      string
	PollInterval      time.Duration
	ProjectList       []string
}

// DefaultConfig provides a default configuration to create a new GitHub
// Eventer by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,
	}
}

// New creates a new configured GitHub Eventer.
func New(config Config) (spec.Eventer, error) {
	if config.Environment == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "environment must not be empty")
	}
	if config.HTTPClientTimeout.Seconds() == 0 {
		return nil, microerror.MaskAnyf(invalidConfigError, "http client timeout must be greater than zero")
	}
	if config.OAuthToken == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "oauth token must not be empty")
	}
	if config.Organisation == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "organisation must not be empty")
	}
	if config.PollInterval.Seconds() == 0 {
		return nil, microerror.MaskAnyf(invalidConfigError, "interval must be greater than zero")
	}
	if len(config.ProjectList) == 0 {
		return nil, microerror.MaskAnyf(invalidConfigError, "project list must not be empty")
	}

	eventer := &githubEventer{
		// Dependencies.
		client: &http.Client{
			Timeout: config.HTTPClientTimeout,
		},
		logger: config.Logger,

		// Settings.
		environment:  config.Environment,
		oauthToken:   config.OAuthToken,
		organisation: config.Organisation,
		pollInterval: config.PollInterval,
		projectList:  config.ProjectList,
	}

	return eventer, nil
}

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

func (e *githubEventer) NewDeploymentEvents() (<-chan spec.DeploymentEvent, error) {
	e.logger.Log("debug", "starting polling for github deployment events", "interval", e.pollInterval)

	deploymentEventChannel := make(chan spec.DeploymentEvent)
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

func (e *githubEventer) SetPending(event spec.DeploymentEvent) error {
	return e.postDeploymentStatus(event.Name, event.ID, pendingState)
}

func (e *githubEventer) SetSuccess(event spec.DeploymentEvent) error {
	return e.postDeploymentStatus(event.Name, event.ID, successState)
}

func (e *githubEventer) SetFailed(event spec.DeploymentEvent) error {
	return e.postDeploymentStatus(event.Name, event.ID, failedState)
}
