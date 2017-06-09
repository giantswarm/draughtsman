package configurer

import (
	"github.com/spf13/viper"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/service/deployer/installer/configurer/configmap"
	"github.com/giantswarm/draughtsman/service/deployer/installer/configurer/file"
	"github.com/giantswarm/draughtsman/service/deployer/installer/configurer/spec"
)

// Config represents the configuration used to create a Configurer.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper

	Type spec.ConfigurerType
}

// DefaultConfig provides a default configuration to create a new Configurer
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Flag:  nil,
		Viper: nil,
	}
}

// New creates a new configured Configurer.
func New(config Config) (spec.Configurer, error) {
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

	var newConfigurer spec.Configurer

	switch config.Type {
	case configmap.ConfigmapConfigurerType:
		configmapConfig := configmap.DefaultConfig()

		configmapConfig.Logger = config.Logger

		configmapConfig.Address = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.Configmap.Address)
		configmapConfig.CAFilePath = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.Configmap.CAFilePath)
		configmapConfig.CrtFilePath = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.Configmap.CrtFilePath)
		configmapConfig.InCluster = config.Viper.GetBool(config.Flag.Service.Deployer.Installer.Configurer.Configmap.InCluster)
		configmapConfig.Key = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.Configmap.Key)
		configmapConfig.KeyFilePath = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.Configmap.KeyFilePath)
		configmapConfig.Name = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.Configmap.Name)
		configmapConfig.Namespace = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.Configmap.Namespace)

		newConfigurer, err = configmap.New(configmapConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

	case file.FileConfigurerType:
		fileConfig := file.DefaultConfig()

		fileConfig.Logger = config.Logger

		fileConfig.Path = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.File.Path)

		newConfigurer, err = file.New(fileConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

	default:
		return nil, microerror.MaskAnyf(invalidConfigError, "configurer type not implemented")
	}

	return newConfigurer, nil
}
