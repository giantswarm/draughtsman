package helm

import (
	micrologger "github.com/giantswarm/microkit/logger"

	eventerspec "github.com/giantswarm/draughtsman/service/deployer/eventer/spec"
	"github.com/giantswarm/draughtsman/service/deployer/installer/spec"
)

// HelmInstallerType is an Installer that uses Helm.
var HelmInstallerType spec.InstallerType = "HelmInstaller"

// Config represents the configuration used to create a Helm Installer.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger
}

// DefaultConfig provides a default configuration to create a new Helm
// Installer by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,
	}
}

// New creates a new configured Helm Installer.
func New(config Config) (*HelmInstaller, error) {
	installer := &HelmInstaller{
		// Dependencies.
		logger: config.Logger,
	}

	return installer, nil
}

// HelmInstaller is an implementation of the Helm interface,
// that uses Helm to install charts.
type HelmInstaller struct {
	// Dependencies.
	logger micrologger.Logger
}

func (i *HelmInstaller) Install(event eventerspec.DeploymentEvent) error {
	i.logger.Log("debug", "installing chart", "name", event.Name)

	return nil
}
