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

	configurerspec "github.com/giantswarm/draughtsman/service/configurer/spec"
	eventerspec "github.com/giantswarm/draughtsman/service/eventer/spec"
	"github.com/giantswarm/draughtsman/service/installer/spec"
)

const (
	// versionedChartFormat is the format the CNR registry uses to address
	// charts. For example, we use this to address that chart to pull.
	versionedChartFormat = "%v/%v/%v-chart@1.0.0-%v"

	// tarballNameFormat is the format for the name of the chart tarball.
	tarballNameFormat = "%v_%v-chart_1.0.0-%v.tar.gz"
)

// HelmInstallerType is an Installer that uses Helm.
var HelmInstallerType spec.InstallerType = "HelmInstaller"

// Config represents the configuration used to create a Helm Installer.
type Config struct {
	// Dependencies.
	Configurers []configurerspec.Configurer
	Logger      micrologger.Logger

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
		Configurers: nil,
		Logger:      nil,
	}
}

// New creates a new configured Helm Installer.
func New(config Config) (*HelmInstaller, error) {
	if config.Configurers == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "configurers must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

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
		configurers: config.Configurers,
		logger:      config.Logger,

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

// HelmInstaller is an implementation of the Installer interface,
// that uses Helm to install charts.
type HelmInstaller struct {
	// Dependencies.
	configurers []configurerspec.Configurer
	logger      micrologger.Logger

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
		versionedChartFormat,
		i.registry,
		i.organisation,
		project,
		sha,
	)
}

// tarballName builds a tarball name, given a project name and sha.
func (i *HelmInstaller) tarballName(project, sha string) string {
	return fmt.Sprintf(
		tarballNameFormat,
		i.organisation,
		project,
		sha,
	)
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
		return microerror.MaskAnyf(err, stdErrBuf.String())
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

	dir, err := os.Getwd()
	if err != nil {
		return microerror.MaskAny(err)
	}

	tarballPath := path.Join(dir, i.tarballName(project, sha))
	if _, err := os.Stat(tarballPath); os.IsNotExist(err) {
		return microerror.MaskAnyf(helmError, "could not find downloaded tarball")
	}

	defer os.Remove(tarballPath)

	i.logger.Log("debug", "downloaded chart", "tarball", tarballPath)

	// The intaller accepts multiple configurers during initialization. Here we
	// iterate over all of them to get all the files they provide.
	var valuesFilesArgs []string
	for _, c := range i.configurers {
		fileName, err := c.File()
		if err != nil {
			return microerror.MaskAny(err)
		}
		valuesFilesArgs = append(valuesFilesArgs, "--values", fileName)
	}

	// The arguments used to execute Helm for app installation can take multiple
	// values files. At the end the command looks something like this.
	//
	//     helm upgrade --install --values ${file1} --values $(file2) ${project} ${tarball_path}
	//
	var installCommand []string
	{
		installCommand = append(installCommand, "upgrade", "--install")
		installCommand = append(installCommand, valuesFilesArgs...)
		installCommand = append(installCommand, project, tarballPath)

		err := i.runHelmCommand("install", installCommand...)
		if err != nil {
			return microerror.MaskAny(err)
		}
	}

	return nil
}
