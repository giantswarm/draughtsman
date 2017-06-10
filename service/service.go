// Package service implements business logic to create Kubernetes resources
// against the Kubernetes API.
package service

import (
	"github.com/spf13/viper"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/http"
	"github.com/giantswarm/draughtsman/service/deployer"
	"github.com/giantswarm/draughtsman/service/version"
)

// Config represents the configuration used to create a new service.
type Config struct {
	// Dependencies.
	HTTPClient http.Client
	Logger     micrologger.Logger

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
		HTTPClient: nil,
		Logger:     nil,

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
		return nil, microerror.MaskAnyf(invalidConfigError, "flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "viper must not be empty")
	}

	var err error

	var deployerService deployer.Deployer
	{
		deployerConfig := deployer.DefaultConfig()

		deployerConfig.HTTPClient = config.HTTPClient
		deployerConfig.Logger = config.Logger

		deployerConfig.Flag = config.Flag
		deployerConfig.Viper = config.Viper

		deployerConfig.Type = deployer.DeployerType(config.Viper.GetString(config.Flag.Service.Deployer.Type))

		deployerService, err = deployer.New(deployerConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
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
			return nil, microerror.MaskAny(err)
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
