package spec

import (
	"github.com/giantswarm/draughtsman/service/eventer/spec"
)

// InstallerType represents the type of Installer to configure.
type InstallerType string

// Installer represents a Service that installs charts.
type Installer interface {
	// Install takes a DeploymentEvent, and installs the referenced chart.
	// If an error occurs, the returned error will be non-nil.
	Install(spec.DeploymentEvent) error
}
