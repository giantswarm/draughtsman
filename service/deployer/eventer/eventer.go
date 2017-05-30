package eventer

import (
	"time"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/draughtsman/service/deployer/eventer/github"
	"github.com/giantswarm/draughtsman/service/deployer/eventer/spec"
)

// Config represents the configuration used to create an Eventer.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Type spec.EventerType

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
		Type: github.GithubEventer,
	}
}

// New creates a new configured Eventer.
func New(config Config) (spec.Eventer, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	var err error

	var newEventer spec.Eventer

	switch config.Type {
	case github.GithubEventer:
		githubConfig := github.DefaultConfig()

		githubConfig.Logger = config.Logger

		githubConfig.Environment = config.Environment
		githubConfig.HTTPClientTimeout = config.HTTPClientTimeout
		githubConfig.OauthToken = config.OauthToken
		githubConfig.Organisation = config.Organisation
		githubConfig.PollInterval = config.PollInterval
		githubConfig.ProjectList = config.ProjectList

		newEventer, err = github.New(githubConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

	default:
		return nil, microerror.MaskAnyf(invalidConfigError, "could not find eventer type")
	}

	return newEventer, nil
}
