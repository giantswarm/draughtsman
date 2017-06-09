package spec

import (
	"github.com/giantswarm/draughtsman/service/deployer/eventer/spec"
)

// NotifierType represents the type of Notifier to configure.
type NotifierType string

// Notifier represents a Service that notifies of install status.
type Notifier interface {
	// Success notifies of successful installations.
	Success(spec.DeploymentEvent) error

	// Failed notifies of failed installations.
	Failed(spec.DeploymentEvent, string) error
}
