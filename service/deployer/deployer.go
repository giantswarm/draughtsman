package deployer

import (
	"github.com/spf13/viper"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/service/deployer/eventer"
	eventerspec "github.com/giantswarm/draughtsman/service/deployer/eventer/spec"
	"github.com/giantswarm/draughtsman/service/deployer/installer"
	installerspec "github.com/giantswarm/draughtsman/service/deployer/installer/spec"
)

// DeployerType represents the type of Deployer to configure.
type DeployerType string

// Config represents the configuration used to create a Deployer.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper

	Type DeployerType
}

// DefaultConfig provides a default configuration to create a new Deployer
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Flag:  nil,
		Viper: nil,

		Type: StandardDeployer,
	}
}

// New creates a new configured Deployer.
func New(config Config) (Deployer, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	var err error

	var eventerService eventerspec.Eventer
	{
		eventerConfig := eventer.DefaultConfig()

		eventerConfig.Logger = config.Logger

		eventerConfig.Flag = config.Flag
		eventerConfig.Viper = config.Viper

		eventerService, err = eventer.New(eventerConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var installerService installerspec.Installer
	{
		installerConfig := installer.DefaultConfig()

		installerConfig.Logger = config.Logger

		installerConfig.Flag = config.Flag
		installerConfig.Viper = config.Viper

		installerService, err = installer.New(installerConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var newService Deployer

	switch config.Type {
	case StandardDeployer:
		newService = &standardDeployer{
			// Dependencies.
			logger:    config.Logger,
			eventer:   eventerService,
			installer: installerService,
		}
	default:
		return nil, microerror.MaskAnyf(invalidConfigError, "could not find deployer type")
	}

	return newService, nil
}

var StandardDeployer DeployerType = "StandardDeployer"

// standardDeployer is an implementation of the Deployer interface.
type standardDeployer struct {
	// Dependencies.
	logger    micrologger.Logger
	eventer   eventerspec.Eventer
	installer installerspec.Installer
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
		} else {
			if err := s.eventer.SetFailed(deploymentEvent); err != nil {
				s.logger.Log("error", "could not set failed event", "message", err.Error())
			}
		}
	}

	s.logger.Log("debug", "finished deployment loop")
}
