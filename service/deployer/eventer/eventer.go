package eventer

import (
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
}

// DefaultConfig provides a default configuration to create a new Eventer
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Type: StubEventer,
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
	case StubEventer:
		newService = &stubEventer{
			// Dependencies.
			logger: config.Logger,
		}
	default:
		return nil, microerror.MaskAnyf(invalidConfigError, "could not find eventer type")
	}

	return newService, nil
}

// StubEventer is an Eventer that just pretends it's a real Eventer.
var StubEventer EventerType = "StubEventer"

// stubEventer is a stub implementation of the Eventer interface.
type stubEventer struct {
	// Dependencies.
	logger micrologger.Logger
}

func (e *stubEventer) NewDeploymentEvents() ([]DeploymentEvent, time.Duration, error) {
	e.logger.Log("debug", "checking for deployment requests")

	event := DeploymentEvent{
		Name: "test-project",
	}

	return []DeploymentEvent{event}, 10 * time.Second, nil
}

func (e *stubEventer) SetPending(event DeploymentEvent) error {
	e.logger.Log("debug", "setting pending", "event name", event.Name)

	return nil
}

func (e *stubEventer) SetSuccess(event DeploymentEvent) error {
	e.logger.Log("debug", "setting success", "event name", event.Name)

	return nil
}

func (e *stubEventer) SetFailed(event DeploymentEvent) error {
	e.logger.Log("debug", "setting failed", "event name", event.Name)

	return nil
}
