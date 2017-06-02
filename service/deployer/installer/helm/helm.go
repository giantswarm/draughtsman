package helm

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	eventerspec "github.com/giantswarm/draughtsman/service/deployer/eventer/spec"
	"github.com/giantswarm/draughtsman/service/deployer/installer/spec"
)

// HelmInstallerType is an Installer that uses Helm.
var HelmInstallerType spec.InstallerType = "HelmInstaller"

// Config represents the configuration used to create a Helm Installer.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	HelmImage    string
	HelmImageTag string
	Password     string
	Registry     string
	Username     string
}

// DefaultConfig provides a default configuration to create a new Helm
// Installer by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,
	}
}

// New creates a new configured Helm Installer.
func New(config Config) (*HelmInstaller, error) {
	if config.HelmImage == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "helm image must not be empty")
	}
	if config.HelmImageTag == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "helm image tag must not be empty")
	}
	if config.Password == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "password must not be empty")
	}
	if config.Registry == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "registry must not be empty")
	}
	if config.Username == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "username must not be empty")
	}

	installer := &HelmInstaller{
		// Dependencies.
		logger: config.Logger,

		// Settings.
		helmImage:    config.HelmImage,
		helmImageTag: config.HelmImageTag,
		password:     config.Password,
		registry:     config.Registry,
		username:     config.Username,
	}

	if err := installer.login(); err != nil {
		return nil, microerror.MaskAny(err)
	}

	return installer, nil
}

// HelmInstaller is an implementation of the Helm interface,
// that uses Helm to install charts.
type HelmInstaller struct {
	// Dependencies.
	logger micrologger.Logger

	// Settings.
	helmImage    string
	helmImageTag string
	password     string
	registry     string
	username     string
}

// dockerImage returns a correctly formatted image string to use with Docker.
func (i *HelmInstaller) dockerImage() string {
	return fmt.Sprintf("%v:%v", i.helmImage, i.helmImageTag)
}

// runHelmCommand runs the given commands inside a Helm container.
func (i *HelmInstaller) runHelmCommand(name string, args ...string) error {
	i.logger.Log("debug", "running helm command", "name", name)

	startTime := time.Now()

	dockerArgs := []string{"run", i.dockerImage()}
	dockerArgs = append(dockerArgs, args...)

	cmd := exec.Command("docker", dockerArgs...)

	var stdOutBuf, stdErrBuf bytes.Buffer
	cmd.Stdout = &stdOutBuf
	cmd.Stderr = &stdErrBuf

	err := cmd.Run()

	i.logger.Log(
		"debug", "ran helm command", "name", name,
		"stdout", stdOutBuf.String(), "stderr", stdErrBuf.String(),
	)

	if err != nil {
		return microerror.MaskAny(err)
	}

	updateHelmMetrics(name, startTime)

	return nil
}

// login logs the configured user into the configured registry.
func (i *HelmInstaller) login() error {
	i.logger.Log("debug", "logging into registry", "username", i.username, "registry", i.registry)

	return i.runHelmCommand(
		"login",
		"registry",
		"login",
		fmt.Sprintf("--user=%v", i.username),
		fmt.Sprintf("--password=%v", i.password),
		i.registry,
	)
}

func (i *HelmInstaller) Install(event eventerspec.DeploymentEvent) error {
	i.logger.Log("debug", "installing chart", "name", event.Name)

	return nil
}
