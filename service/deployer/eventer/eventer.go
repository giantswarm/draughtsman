package eventer

import (
	"github.com/spf13/viper"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/service/deployer/eventer/github"
	"github.com/giantswarm/draughtsman/service/deployer/eventer/spec"
)

// Config represents the configuration used to create an Eventer.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

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
		Logger: nil,

		// Settings.
		Flag:  nil,
		Viper: nil,
	}
}

// New creates a new configured Eventer.
func New(config Config) (spec.Eventer, error) {
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

	var newEventer spec.Eventer

	switch config.Type {
	case github.GithubEventerType:
		githubConfig := github.DefaultConfig()

		githubConfig.Logger = config.Logger

		githubConfig.Environment = config.Viper.GetString(config.Flag.Service.Deployer.Environment)
		githubConfig.HTTPClientTimeout = config.Viper.GetDuration(config.Flag.Service.Deployer.Eventer.GitHub.HTTPClientTimeout)
		githubConfig.OAuthToken = config.Viper.GetString(config.Flag.Service.Deployer.Eventer.GitHub.OAuthToken)
		githubConfig.Organisation = config.Viper.GetString(config.Flag.Service.Deployer.Eventer.GitHub.Organisation)
		githubConfig.PollInterval = config.Viper.GetDuration(config.Flag.Service.Deployer.Eventer.GitHub.PollInterval)
		githubConfig.ProjectList = config.Viper.GetStringSlice(config.Flag.Service.Deployer.Eventer.GitHub.ProjectList)

		newEventer, err = github.New(githubConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

	default:
		return nil, microerror.MaskAnyf(invalidConfigError, "eventer type not implemented")
	}

	return newEventer, nil
}
