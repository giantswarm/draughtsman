package installer

import (
	"github.com/giantswarm/draughtsman/service/deployer/eventer"
)

// Installer represents a Service that installs packages.
type Installer interface {
	// Install takes a DeploymentEvent, and installs the referenced package.
	// If an error occurs, the returned error will be non-nil.
	Install(eventer.DeploymentEvent) error
}
