package deployer

import (
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/service/eventer"
	eventerspec "github.com/giantswarm/draughtsman/service/eventer/spec"
	httpspec "github.com/giantswarm/draughtsman/service/http"
	"github.com/giantswarm/draughtsman/service/installer"
	installerspec "github.com/giantswarm/draughtsman/service/installer/spec"
	"github.com/giantswarm/draughtsman/service/notifier"
	notifierspec "github.com/giantswarm/draughtsman/service/notifier/spec"
	slackspec "github.com/giantswarm/draughtsman/service/slack"
)

// DeployerType represents the type of Deployer to configure.
type DeployerType string

// Config represents the configuration used to create a Deployer.
type Config struct {
	// Dependencies.
	FileSystem       afero.Fs
	HTTPClient       httpspec.Client
	KubernetesClient kubernetes.Interface
	Logger           micrologger.Logger
	SlackClient      slackspec.Client

	// Settings.
	Flag  *flag.Flag
	Type  DeployerType
	Viper *viper.Viper
}

// DefaultConfig provides a default configuration to create a new Deployer
// service by best effort.
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
		Type:  "",
		Viper: nil,
	}
}

// New creates a new configured Deployer.
func New(config Config) (Deployer, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "logger must not be empty")
	}

	// Settings.
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "viper must not be empty")
	}

	var err error

	var eventerService eventerspec.Eventer
	{
		eventerConfig := eventer.DefaultConfig()

		eventerConfig.HTTPClient = config.HTTPClient
		eventerConfig.Logger = config.Logger

		eventerConfig.Flag = config.Flag
		eventerConfig.Viper = config.Viper

		eventerConfig.Type = eventerspec.EventerType(config.Viper.GetString(config.Flag.Service.Deployer.Eventer.Type))

		eventerService, err = eventer.New(eventerConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var installerService installerspec.Installer
	{
		installerConfig := installer.DefaultConfig()

		installerConfig.FileSystem = config.FileSystem
		installerConfig.KubernetesClient = config.KubernetesClient
		installerConfig.Logger = config.Logger

		installerConfig.Flag = config.Flag
		installerConfig.Viper = config.Viper

		installerConfig.Type = installerspec.InstallerType(config.Viper.GetString(config.Flag.Service.Deployer.Installer.Type))

		installerService, err = installer.New(installerConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var notifierService notifierspec.Notifier
	{
		notifierConfig := notifier.DefaultConfig()

		notifierConfig.Logger = config.Logger
		notifierConfig.SlackClient = config.SlackClient

		notifierConfig.Flag = config.Flag
		notifierConfig.Viper = config.Viper

		notifierConfig.Type = notifierspec.NotifierType(config.Viper.GetString(config.Flag.Service.Deployer.Notifier.Type))

		notifierService, err = notifier.New(notifierConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var newService Deployer
	switch config.Type {
	case StandardDeployer:
		newService = &standardDeployer{
			// Dependencies.
			eventer:   eventerService,
			installer: installerService,
			logger:    config.Logger,
			notifier:  notifierService,
		}
	default:
		return nil, microerror.Maskf(invalidConfigError, "could not find deployer type")
	}

	return newService, nil
}

var StandardDeployer DeployerType = "StandardDeployer"

// standardDeployer is an implementation of the Deployer interface.
type standardDeployer struct {
	// Dependencies.
	eventer   eventerspec.Eventer
	installer installerspec.Installer
	logger    micrologger.Logger
	notifier  notifierspec.Notifier
}

// Boot starts the deployer.
func (s *standardDeployer) Boot() {
	s.logger.Log("debug", "starting deployer")

	deploymentEventChannel, err := s.eventer.NewDeploymentEvents()
	if err != nil {
		s.logger.Log("debug", "could not get deployment event channel", "message", err.Error())
	}

	for deploymentEvent := range deploymentEventChannel {
		if err := s.eventer.SetPending(deploymentEvent); err != nil {
			s.logger.Log("error", "could not set pending event", "message", err.Error())
		}

		installErr := s.installer.Install(deploymentEvent)
		if installErr == nil {
			if err := s.eventer.SetSuccess(deploymentEvent); err != nil {
				s.logger.Log("error", "could not set success event", "message", err.Error())
			}

			if err := s.notifier.Success(deploymentEvent); err != nil {
				s.logger.Log("error", "could not notify of success", "message", err.Error())
			}
		} else {
			s.logger.Log("error", "could not install chart", "message", installErr.Error())

			if err := s.eventer.SetFailed(deploymentEvent); err != nil {
				s.logger.Log("error", "could not set failed event", "message", err.Error())
			}

			if err := s.notifier.Failed(deploymentEvent, installErr.Error()); err != nil {
				s.logger.Log("error", "could not notify of failure", "message", err.Error())
			}
		}
	}

	s.logger.Log("debug", "finished deployment loop")
}
