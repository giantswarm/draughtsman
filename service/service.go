// Package service implements business logic to create Kubernetes resources
// against the Kubernetes API.
package service

import (
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/service/deployer"
	httpspec "github.com/giantswarm/draughtsman/service/http"
	slackspec "github.com/giantswarm/draughtsman/service/slack"
	"github.com/giantswarm/draughtsman/service/version"
)

// Config represents the configuration used to create a new service.
type Config struct {
	// Dependencies.
	FileSystem       afero.Fs
	HTTPClient       httpspec.Client
	KubernetesClient kubernetes.Interface
	Logger           micrologger.Logger
	SlackClient      slackspec.Client

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper

	Description string
	GitCommit   string
	Name        string
	Source      string
}

// DefaultConfig provides a default configuration to create a new service by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		FileSystem:       afero.NewMemMapFs(),
		HTTPClient:       nil,
		KubernetesClient: nil,
		Logger:           nil,
		SlackClient:      nil,

		// Settings.
		Flag:  nil,
		Viper: nil,

		Description: "",
		GitCommit:   "",
		Name:        "",
		Source:      "",
	}
}

// New creates a new configured service object.
func New(config Config) (*Service, error) {
	// Settings.
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "viper must not be empty")
	}

	var err error

	var deployerService deployer.Deployer
	{
		deployerConfig := deployer.DefaultConfig()

		deployerConfig.FileSystem = config.FileSystem
		deployerConfig.HTTPClient = config.HTTPClient
		deployerConfig.KubernetesClient = config.KubernetesClient
		deployerConfig.Logger = config.Logger
		deployerConfig.SlackClient = config.SlackClient

		deployerConfig.Flag = config.Flag
		deployerConfig.Viper = config.Viper

		deployerConfig.Type = deployer.DeployerType(config.Viper.GetString(config.Flag.Service.Deployer.Type))

		deployerService, err = deployer.New(deployerConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var versionService *version.Service
	{
		versionConfig := version.DefaultConfig()

		versionConfig.Description = config.Description
		versionConfig.GitCommit = config.GitCommit
		versionConfig.Name = config.Name
		versionConfig.Source = config.Source

		versionService, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newService := &Service{
		// Dependencies.
		Deployer: deployerService,
		Version:  versionService,
	}

	return newService, nil
}

type Service struct {
	// Dependencies.
	Deployer deployer.Deployer
	Version  *version.Service
}

func (s *Service) Boot() {
	s.Deployer.Boot()
}
