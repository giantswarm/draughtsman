package github

import (
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/draughtsman/service/eventer/spec"
	httpspec "github.com/giantswarm/draughtsman/service/http"
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
		return nil, microerror.Maskf(invalidConfigError, "http client must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "logger must not be empty")
	}

	if config.Environment == "" {
		return nil, microerror.Maskf(invalidConfigError, "environment must not be empty")
	}
	if config.OAuthToken == "" {
		return nil, microerror.Maskf(invalidConfigError, "oauth token must not be empty")
	}
	if config.Organisation == "" {
		return nil, microerror.Maskf(invalidConfigError, "organisation must not be empty")
	}
	if config.PollInterval.Seconds() == 0 {
		return nil, microerror.Maskf(invalidConfigError, "interval must be greater than zero")
	}
	if len(config.ProjectList) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "project list must not be empty")
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
			e.logger.Log("debug", "Fetching deployment events", "projectlist", e.projectList)
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
