package installer

import (
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/microerror"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/service/configurer"
	configurerspec "github.com/giantswarm/draughtsman/service/configurer/spec"
	"github.com/giantswarm/draughtsman/service/installer/helm"
	"github.com/giantswarm/draughtsman/service/installer/spec"
)

// Config represents the configuration used to create an Installer.
type Config struct {
	// Dependencies.
	FileSystem       afero.Fs
	KubernetesClient kubernetes.Interface
	Logger           micrologger.Logger

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
		FileSystem:       afero.NewMemMapFs(),
		KubernetesClient: nil,
		Logger:           nil,

		// Settings.
		Flag:  nil,
		Viper: nil,
	}
}

// New creates a new configured Installer.
func New(config Config) (spec.Installer, error) {
	// Settings.
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "viper must not be empty")
	}

	var err error

	var configurerServices []configurerspec.Configurer
	types := strings.Split(config.Viper.GetString(config.Flag.Service.Deployer.Installer.Configurer.Types), ",")
	for _, t := range types {
		configurerConfig := configurer.DefaultConfig()

		configurerConfig.FileSystem = config.FileSystem
		configurerConfig.KubernetesClient = config.KubernetesClient
		configurerConfig.Logger = config.Logger

		configurerConfig.Flag = config.Flag
		configurerConfig.Type = configurerspec.ConfigurerType(t)
		configurerConfig.Viper = config.Viper

		configurerService, err := configurer.New(configurerConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		configurerServices = append(configurerServices, configurerService)
	}

	var newInstaller spec.Installer
	switch config.Type {
	case helm.HelmInstallerType:
		helmConfig := helm.DefaultConfig()

		helmConfig.Configurers = configurerServices
		helmConfig.FileSystem = config.FileSystem
		helmConfig.Logger = config.Logger

		helmConfig.HelmBinaryPath = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Helm.HelmBinaryPath)
		helmConfig.Organisation = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Helm.Organisation)
		helmConfig.Password = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Helm.Password)
		helmConfig.Registry = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Helm.Registry)
		helmConfig.Username = config.Viper.GetString(config.Flag.Service.Deployer.Installer.Helm.Username)

		newInstaller, err = helm.New(helmConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}

	default:
		return nil, microerror.Maskf(invalidConfigError, "installer type not implemented")
	}

	return newInstaller, nil
}
