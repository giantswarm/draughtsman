package helm

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
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
	HelmBinaryPath string
	Organisation   string
	Password       string
	Registry       string
	Username       string
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
	if config.HelmBinaryPath == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "helm binary path must not be empty")
	}
	if config.Organisation == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "organisation must not be empty")
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

	if _, err := os.Stat(config.HelmBinaryPath); os.IsNotExist(err) {
		return nil, microerror.MaskAnyf(invalidConfigError, "helm binary does not exist")
	}

	installer := &HelmInstaller{
		// Dependencies.
		logger: config.Logger,

		// Settings.
		helmBinaryPath: config.HelmBinaryPath,
		organisation:   config.Organisation,
		password:       config.Password,
		registry:       config.Registry,
		username:       config.Username,
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
	helmBinaryPath string
	organisation   string
	password       string
	registry       string
	username       string
}

// versionedChartName builds a chart name, including a version,
// given a project name and a sha.
func (i *HelmInstaller) versionedChartName(project, sha string) string {
	return fmt.Sprintf(
		"%v/%v/%v-chart@1.0.0-%v",
		i.registry,
		i.organisation,
		project,
		sha,
	)
}

// tarballPath creates a tarball path, given a project name and a sha.
func (i *HelmInstaller) tarballPath(project, sha string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", microerror.MaskAny(err)
	}

	tarballName := fmt.Sprintf(
		"%v_%v-chart_1.0.0-%v.tar.gz",
		i.organisation,
		project,
		sha,
	)

	return path.Join(dir, tarballName), nil
}

// runHelmCommand runs the given Helm command.
func (i *HelmInstaller) runHelmCommand(name string, args ...string) error {
	i.logger.Log("debug", "running helm command", "name", name)

	startTime := time.Now()

	cmd := exec.Command(i.helmBinaryPath, args...)

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

	if strings.Contains(stdOutBuf.String(), "Error") {
		return microerror.MaskAnyf(helmError, stdOutBuf.String())
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
	project := event.Name
	sha := event.Sha

	i.logger.Log("debug", "installing chart", "name", project, "sha", sha)

	if err := i.runHelmCommand(
		"pull",
		"registry",
		"pull",
		i.versionedChartName(project, sha),
	); err != nil {
		return microerror.MaskAny(err)
	}

	tarballPath, err := i.tarballPath(project, sha)
	if err != nil {
		return microerror.MaskAny(err)
	}
	if _, err := os.Stat(tarballPath); os.IsNotExist(err) {
		return microerror.MaskAnyf(helmError, "could not find downloaded tarball")
	}
	defer os.Remove(tarballPath)

	i.logger.Log("debug", "downloaded chart", "tarball", tarballPath)

	if err := i.runHelmCommand(
		"install",
		"install",
		tarballPath,
		// "--values",
		// "./values.yaml",
		"--wait",
	); err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}
