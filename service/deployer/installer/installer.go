package installer

import (
	"github.com/spf13/viper"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/service/deployer/installer/helm"
	"github.com/giantswarm/draughtsman/service/deployer/installer/spec"
)

// Config represents the configuration used to create an Installer.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper

	Type spec.InstallerType
}

// DefaultConfig provides a default configuration to create a new Installer
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Flag:  nil,
		Viper: nil,

		Type: helm.HelmInstallerType,
	}
}

// New creates a new configured Installer.
func New(config Config) (spec.Installer, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	// Settings.
	if config.Flag == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "viper must not be empty")
	}

	var err error

	var newInstaller spec.Installer

	switch config.Type {
	case helm.HelmInstallerType:
		helmConfig := helm.DefaultConfig()

		helmConfig.Logger = config.Logger

		helmConfig.HelmBinaryPath = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Helm.HelmBinaryPath)
		helmConfig.Organisation = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Helm.Organisation)
		helmConfig.Password = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Helm.Password)
		helmConfig.Registry = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Helm.Registry)
		helmConfig.Username = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Helm.Username)

		newInstaller, err = helm.New(helmConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

	default:
		return nil, microerror.MaskAnyf(invalidConfigError, "installer type not implemented")
	}

	return newInstaller, nil
}
