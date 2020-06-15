package eventer

import (
	"github.com/spf13/viper"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/pkg/project"
	"github.com/giantswarm/draughtsman/service/eventer/github"
	"github.com/giantswarm/draughtsman/service/eventer/spec"
	httpspec "github.com/giantswarm/draughtsman/service/http"
)

// Config represents the configuration used to create an Eventer.
type Config struct {
	// Dependencies.
	HTTPClient httpspec.Client
	Logger     micrologger.Logger

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper

	Type spec.EventerType
}

// DefaultConfig provides a default configuration to create a new Eventer
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		HTTPClient: nil,
		Logger:     nil,

		// Settings.
		Flag:  nil,
		Viper: nil,
	}
}

// New creates a new configured Eventer.
func New(config Config) (spec.Eventer, error) {
	// Settings.
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "viper must not be empty")
	}

	var err error

	var newEventer spec.Eventer
	switch config.Type {
	case github.GithubEventerType:
		githubConfig := github.DefaultConfig()

		githubConfig.HTTPClient = config.HTTPClient
		githubConfig.Logger = config.Logger

		githubConfig.Environment = config.Viper.GetString(config.Flag.Service.Deployer.Environment)
		githubConfig.OAuthToken = config.Viper.GetString(config.Flag.Service.Deployer.Eventer.GitHub.OAuthToken)
		githubConfig.Organisation = config.Viper.GetString(config.Flag.Service.Deployer.Eventer.GitHub.Organisation)
		githubConfig.PollInterval = config.Viper.GetDuration(config.Flag.Service.Deployer.Eventer.GitHub.PollInterval)
		githubConfig.Provider = config.Viper.GetString(config.Flag.Service.Deployer.Provider)

		{
			projectList := project.GetProjectList(githubConfig.Provider, githubConfig.Environment)
			githubConfig.ProjectList = projectList
		}

		newEventer, err = github.New(githubConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}

	default:
		return nil, microerror.Maskf(invalidConfigError, "eventer type not implemented")
	}

	return newEventer, nil
}
