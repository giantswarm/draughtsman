package notifier

import (
	"github.com/spf13/viper"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/service/notifier/slack"
	"github.com/giantswarm/draughtsman/service/notifier/spec"
	slackspec "github.com/giantswarm/draughtsman/service/slack"
)

// Config represents the configuration used to create a Notifier.
type Config struct {
	// Dependencies.
	Logger      micrologger.Logger
	SlackClient slackspec.Client

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper

	Type spec.NotifierType
}

// DefaultConfig provides a default configuration to create a new Notifier
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:      nil,
		SlackClient: nil,

		// Settings.
		Flag:  nil,
		Viper: nil,
	}
}

// New creates a new configured Notifier.
func New(config Config) (spec.Notifier, error) {
	// Settings.
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "viper must not be empty")
	}

	var err error

	var newNotifier spec.Notifier
	switch config.Type {
	case slack.SlackNotifierType:
		slackConfig := slack.DefaultConfig()

		slackConfig.Logger = config.Logger
		slackConfig.SlackClient = config.SlackClient

		slackConfig.Channel = config.Viper.GetString(config.Flag.Service.Deployer.Notifier.Slack.Channel)
		slackConfig.Emoji = config.Viper.GetString(config.Flag.Service.Deployer.Notifier.Slack.Emoji)
		slackConfig.Environment = config.Viper.GetString(config.Flag.Service.Deployer.Environment)
		slackConfig.Username = config.Viper.GetString(config.Flag.Service.Deployer.Notifier.Slack.Username)

		newNotifier, err = slack.New(slackConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}

	default:
		return nil, microerror.Maskf(invalidConfigError, "notifier type not implemented")
	}

	return newNotifier, nil
}
