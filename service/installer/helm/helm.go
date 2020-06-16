package helm

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"

	"github.com/giantswarm/draughtsman/pkg/project/configuration"
	configurerspec "github.com/giantswarm/draughtsman/service/configurer/spec"
	eventerspec "github.com/giantswarm/draughtsman/service/eventer/spec"
	"github.com/giantswarm/draughtsman/service/installer/spec"
)

const (
	// versionedChartFormat is the format the CNR registry uses to address
	// charts. For example, we use this to address that chart to pull.
	versionedChartFormat = "%v/%v/%v-chart@1.0.0-%v"

	// chartNameFormat is the format for the name of the chart folder.
	chartNameFormat = "%v_%v-chart_1.0.0-%v/%v-chart"

	// pollInterval is the time interval between checking the status of helm releases
	pollInterval = 1 * time.Minute
)

// HelmInstallerType is an Installer that uses Helm.
var HelmInstallerType spec.InstallerType = "HelmInstaller"

// Config represents the configuration used to create a Helm Installer.
type Config struct {
	// Dependencies.
	Configurers []configurerspec.Configurer
	FileSystem  afero.Fs
	Logger      micrologger.Logger

	// Settings.
	Environment    string
	HelmBinaryPath string
	Organisation   string
	Password       string
	Provider       string
	Registry       string
	Username       string
}

// DefaultConfig provides a default configuration to create a new Helm
// Installer by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Configurers: nil,
		FileSystem:  afero.NewMemMapFs(),
		Logger:      nil,

		// Settings.
		HelmBinaryPath: "",
		Organisation:   "",
		Password:       "",
		Registry:       "",
		Username:       "",
	}
}

// New creates a new configured Helm Installer.
func New(config Config) (*HelmInstaller, error) {
	// Dependencies.
	if config.Configurers == nil {
		return nil, microerror.Maskf(invalidConfigError, "configurers must not be empty")
	}
	if config.FileSystem == nil {
		return nil, microerror.Maskf(invalidConfigError, "file system must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "logger must not be empty")
	}

	// Settings.
	if config.Environment == "" {
		return nil, microerror.Maskf(invalidConfigError, "environment must not be empty")
	}
	if config.HelmBinaryPath == "" {
		return nil, microerror.Maskf(invalidConfigError, "helm binary path must not be empty")
	}
	if config.Organisation == "" {
		return nil, microerror.Maskf(invalidConfigError, "organisation must not be empty")
	}
	if config.Password == "" {
		return nil, microerror.Maskf(invalidConfigError, "password must not be empty")
	}
	if config.Provider == "" {
		return nil, microerror.Maskf(invalidConfigError, "provider must not be empty")
	}

	if config.Registry == "" {
		return nil, microerror.Maskf(invalidConfigError, "registry must not be empty")
	}
	if config.Username == "" {
		return nil, microerror.Maskf(invalidConfigError, "username must not be empty")
	}

	if _, err := os.Stat(config.HelmBinaryPath); os.IsNotExist(err) {
		return nil, microerror.Maskf(invalidConfigError, "helm binary does not exist")
	}

	installer := &HelmInstaller{
		// Dependencies.
		configurers: config.Configurers,
		fileSystem:  config.FileSystem,
		logger:      config.Logger,

		// Settings.
		environment:    config.Environment,
		helmBinaryPath: config.HelmBinaryPath,
		organisation:   config.Organisation,
		password:       config.Password,
		provider:       config.Provider,
		registry:       config.Registry,
		username:       config.Username,
	}

	if err := installer.login(); err != nil {
		return nil, microerror.Mask(err)
	}

	if err := installer.fetchMetrics(); err != nil {
		return nil, microerror.Mask(err)
	}

	return installer, nil
}

// HelmInstaller is an implementation of the Installer interface,
// that uses Helm to install charts.
type HelmInstaller struct {
	// Dependencies.
	configurers []configurerspec.Configurer
	fileSystem  afero.Fs
	logger      micrologger.Logger

	// Settings.
	environment    string
	helmBinaryPath string
	organisation   string
	password       string
	provider       string
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

// chartName builds a chart name, given a project name and sha.
func (i *HelmInstaller) chartName(project, sha string) string {
	return fmt.Sprintf(
		chartNameFormat,
		i.organisation,
		project,
		sha,
		project,
	)
}

// runHelmCommand runs the given Helm command.
func (i *HelmInstaller) runHelmCommand(name string, args ...string) error {
	i.logger.Log("debug", "running helm command", "name", name)

	defer updateHelmMetrics(name, time.Now())

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
		return microerror.Maskf(helmError, "error output: %s", stdErrBuf.String())
	}

	if strings.Contains(stdOutBuf.String(), "Error") {
		return microerror.Maskf(helmError, stdOutBuf.String())
	}

	return nil
}

func (i *HelmInstaller) fetchMetrics() error {
	ticker := time.NewTicker(pollInterval)
	projectList := configuration.GetProjectList(i.provider, i.environment)

	go func() {
		for c := ticker.C; ; <-c {
			i.logger.Log("debug", "fetching metrics")
			i.checkHelmRelease(projectList)
			i.logger.Log("debug", "fetched metrics")
		}
	}()

	return nil
}

func (i *HelmInstaller) checkHelmRelease(projectList []string) {
	for _, prj := range projectList {
		cmd := exec.Command(i.helmBinaryPath, "history", prj, "--output", "yaml", "--max", "1")

		var stdOutBuf, stdErrBuf bytes.Buffer
		cmd.Stdout = &stdOutBuf
		cmd.Stderr = &stdErrBuf

		err := cmd.Run()
		if err != nil {
			i.logger.Log("error", fmt.Sprintf("could not find helm release %#q history", prj), stdErrBuf.String())
		}

		var v []map[string]string
		err = yaml.Unmarshal(stdOutBuf.Bytes(), &v)
		if err != nil {
			i.logger.Log("error", "could not parse helm history output", err.Error())
		}

		helmReleaseFailure.Reset()

		if len(v) > 0 {
			reportHelmRelease(prj, strings.ToLower(v[0]["status"]))
		} else {
			reportHelmRelease(prj, "not-found")
		}
	}
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
		return microerror.Mask(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		return microerror.Mask(err)
	}

	chartPath := path.Join(dir, i.chartName(project, sha))
	if _, err := os.Stat(chartPath); os.IsNotExist(err) {
		return microerror.Maskf(helmError, "could not find downloaded chart")
	}

	defer os.Remove(chartPath)

	i.logger.Log("debug", "downloaded chart", "chart", chartPath)

	// We create a tmp dir in which all Helm values files are written to. After we
	// are done we can just remove the whole tmp dir to clean up.
	var tmpDir string
	{
		tmpDir, err = afero.TempDir(i.fileSystem, "", "draughtsman-installer")
		if err != nil {
			return microerror.Mask(err)
		}
		defer func() {
			err := i.fileSystem.RemoveAll(tmpDir)
			if err != nil {
				i.logger.Log("error", fmt.Sprintf("could not remove tmp dir: %#v", err), "name", project, "sha", sha)
			}
		}()
	}

	// The intaller accepts multiple configurers during initialization. Here we
	// iterate over all of them to get all the values they provide. For each
	// values file we have to create a file in the tmp dir we created above.
	var valuesFilesArgs []string
	for _, c := range i.configurers {
		fileName := filepath.Join(tmpDir, fmt.Sprintf("%s-values.yaml", strings.ToLower(string(c.Type()))))
		values, err := c.Values()
		if err != nil {
			return microerror.Mask(err)
		}

		err = afero.WriteFile(i.fileSystem, fileName, []byte(values), os.FileMode(0644))
		if err != nil {
			return microerror.Mask(err)
		}

		valuesFilesArgs = append(valuesFilesArgs, "--values", fileName)
	}

	// The arguments used to execute Helm for app installation can take multiple
	// values files. At the end the command looks something like this.
	//
	//     helm upgrade --install --values ${file1} --values $(file2) ${project} ${chart_path}
	//
	var installCommand []string
	{
		installCommand = append(installCommand, "upgrade", "--install")
		installCommand = append(installCommand, valuesFilesArgs...)
		installCommand = append(installCommand, project, chartPath)

		err := i.runHelmCommand("install", installCommand...)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
