package installer

import (
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/draughtsman/service/deployer/eventer"
)

// InstallerType represents the type of Installer to configure.
type InstallerType string

// Config represents the configuration used to create an Installer
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Type InstallerType
}

// DefaultConfig provides a default configuration to create a new Installer
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Type: StubInstaller,
	}
}

// New creates a new configured Installer
func New(config Config) (Installer, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	var newService Installer

	switch config.Type {
	case StubInstaller:
		newService = &stubInstaller{
			// Dependencies.
			logger: config.Logger,
		}
	default:
		return nil, microerror.MaskAnyf(invalidConfigError, "could not find installer type")
	}

	return newService, nil
}

// StubInstaller is an Installer that just pretends it's a real Installer.
var StubInstaller InstallerType = "StubInstaller"

// stubInstaller is a stub implementation of the Installer interface.
type stubInstaller struct {
	// Dependencies.
	logger micrologger.Logger
}

func (i *stubInstaller) Install(event eventer.DeploymentEvent) error {
	i.logger.Log("debug", "installing package", "name", event.Name)

	return nil
}
