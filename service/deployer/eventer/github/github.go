package github

import (
	"time"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	httpspec "github.com/giantswarm/draughtsman/http"
	"github.com/giantswarm/draughtsman/service/deployer/eventer/spec"
)

// GithubEventerType is an Eventer that uses Github Deployment Events as a backend.
var GithubEventerType spec.EventerType = "GithubEventer"

// Config represents the configuration used to create a GitHub Eventer.
type Config struct {
	// Dependencies.
	HTTPClient httpspec.Client
	Logger     micrologger.Logger

	Environment  string
	OAuthToken   string
	Organisation string
	PollInterval time.Duration
	ProjectList  []string
}

// DefaultConfig provides a default configuration to create a new GitHub
// Eventer by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		HTTPClient: nil,
		Logger:     nil,
	}
}

// New creates a new configured GitHub Eventer.
func New(config Config) (*GithubEventer, error) {
	if config.HTTPClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "http client must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	if config.Environment == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "environment must not be empty")
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

	eventer := &GithubEventer{
		// Dependencies.
		client: config.HTTPClient,
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

// GithubEventer is an implementation of the Eventer interface,
// that uses GitHub Deployment Events as a backend.
type GithubEventer struct {
	// Dependencies.
	client httpspec.Client
	logger micrologger.Logger

	// Settings.
	environment  string
	oauthToken   string
	organisation string
	pollInterval time.Duration
	projectList  []string
}

func (e *GithubEventer) NewDeploymentEvents() (<-chan spec.DeploymentEvent, error) {
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

func (e *GithubEventer) SetPending(event spec.DeploymentEvent) error {
	return e.postDeploymentStatus(event.Name, event.ID, pendingState)
}

func (e *GithubEventer) SetSuccess(event spec.DeploymentEvent) error {
	return e.postDeploymentStatus(event.Name, event.ID, successState)
}

func (e *GithubEventer) SetFailed(event spec.DeploymentEvent) error {
	return e.postDeploymentStatus(event.Name, event.ID, failureState)
}
