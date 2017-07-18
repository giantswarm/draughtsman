package configurer

import (
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/service/configurer/configmap"
	"github.com/giantswarm/draughtsman/service/configurer/file"
	"github.com/giantswarm/draughtsman/service/configurer/secret"
	"github.com/giantswarm/draughtsman/service/configurer/spec"
)

// Config represents the configuration used to create a Configurer.
type Config struct {
	// Dependencies.
	KubernetesClient kubernetes.Interface
	Logger           micrologger.Logger

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
		KubernetesClient: nil,
		Logger:           nil,

		// Settings.
		Flag:  nil,
		Viper: nil,
	}
}

// New creates a new configured Configurer.
func New(config Config) (spec.Configurer, error) {
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
	case configmap.ConfigMapConfigurerType:
		configmapConfig := configmap.DefaultConfig()

		configmapConfig.KubernetesClient = config.KubernetesClient
		configmapConfig.Logger = config.Logger

		configmapConfig.Key = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.ConfigMap.Key)
		configmapConfig.Name = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.ConfigMap.Name)
		configmapConfig.Namespace = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.ConfigMap.Namespace)

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

	case secret.SecretConfigurerType:
		secretConfig := secret.DefaultConfig()

		secretConfig.KubernetesClient = config.KubernetesClient
		secretConfig.Logger = config.Logger

		secretConfig.Key = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.Secret.Key)
		secretConfig.Name = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.Secret.Name)
		secretConfig.Namespace = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.Secret.Namespace)

		newConfigurer, err = secret.New(secretConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

	default:
		return nil, microerror.MaskAnyf(invalidConfigError, "configurer type not implemented")
	}

	return newConfigurer, nil
}
