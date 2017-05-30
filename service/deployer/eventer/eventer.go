package eventer

import (
	"net/http"
	"time"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
)

// EventerType represents the type of Eventer to configure.
type EventerType string

// Config represents the configuration used to create an Eventer.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Type EventerType

	// GithubEventer settings.
	Environment       string
	HTTPClientTimeout time.Duration
	OauthToken        string
	Organisation      string
	PollInterval      time.Duration
	ProjectList       []string
}

// DefaultConfig provides a default configuration to create a new Eventer
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Type: GithubEventer,
	}
}

// New creates a new configured Eventer.
func New(config Config) (Eventer, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	if config.Environment == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "environment must not be empty")
	}
	if config.HTTPClientTimeout.Seconds() == 0 {
		return nil, microerror.MaskAnyf(invalidConfigError, "http client timeout must be greater than zero")
	}
	if config.OauthToken == "" {
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

	var newService Eventer

	switch config.Type {
	case GithubEventer:
		newService = &githubEventer{
			// Dependencies.
			client: &http.Client{
				Timeout: config.HTTPClientTimeout,
			},
			logger: config.Logger,

			// Settings.
			environment:  config.Environment,
			oauthToken:   config.OauthToken,
			organisation: config.Organisation,
			pollInterval: config.PollInterval,
			projectList:  config.ProjectList,
		}
	default:
		return nil, microerror.MaskAnyf(invalidConfigError, "could not find eventer type")
	}

	return newService, nil
}
