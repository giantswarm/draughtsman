package installer

import (
	"github.com/giantswarm/draughtsman/service/deployer/eventer/spec"
)

// Installer represents a Service that installs charts.
type Installer interface {
	// Install takes a DeploymentEvent, and installs the referenced chart.
	// If an error occurs, the returned error will be non-nil.
	Install(spec.DeploymentEvent) error
}
