package deployer

import (
	"time"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/draughtsman/service/deployer/eventer"
	"github.com/giantswarm/draughtsman/service/deployer/installer"
)

// DeployerType represents the type of Deployer to configure.
type DeployerType string

// Config represents the configuration used to create a Deployer.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Type DeployerType
}

// DefaultConfig provides a default configuration to create a new Deployer
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
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

	var eventerService eventer.Eventer
	{
		eventerConfig := eventer.DefaultConfig()

		eventerConfig.Logger = config.Logger

		eventerService, err = eventer.New(eventerConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var installerService installer.Installer
	{
		installerConfig := installer.DefaultConfig()

		installerConfig.Logger = config.Logger

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
	eventer   eventer.Eventer
	installer installer.Installer
}

// Boot starts the deployer.
func (s *standardDeployer) Boot() {
	s.logger.Log("debug", "starting deployer")

	for {
		deploymentEvents, waitTime, err := s.eventer.NewDeploymentEvents()
		if err != nil {
			s.logger.Log("error", "could not fetch deployment events", "message", err.Error())
		}

		for _, deploymentEvent := range deploymentEvents {
			s.logger.Log("debug", "installing package", "name", deploymentEvent.Name)

			if err := s.eventer.SetPending(deploymentEvent); err != nil {
				s.logger.Log("error", "could not set pending event", "message", err.Error())
			}

			installErr := s.installer.Install(deploymentEvent)
			if installErr == nil {
				s.logger.Log("debug", "successfully installed package", "name", deploymentEvent.Name)

				if err := s.eventer.SetSuccess(deploymentEvent); err != nil {
					s.logger.Log("error", "could not set success event", "message", err.Error())
				}
			} else {
				s.logger.Log("error", "could not install package", "name", deploymentEvent.Name, "message", err.Error())

				if err := s.eventer.SetFailed(deploymentEvent); err != nil {
					s.logger.Log("error", "could not set failed event", "message", err.Error())
				}
			}
		}

		s.logger.Log("debug", "waiting to check deployment events", "wait time", waitTime.String())
		time.Sleep(waitTime)
	}
}
