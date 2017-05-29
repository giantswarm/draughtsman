package eventer

import (
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

	var newService Eventer

	switch config.Type {
	case GithubEventer:
		newService = &githubEventer{
			// Dependencies.
			logger: config.Logger,
		}
	default:
		return nil, microerror.MaskAnyf(invalidConfigError, "could not find eventer type")
	}

	return newService, nil
}
