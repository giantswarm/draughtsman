// Package service implements business logic to create Kubernetes resources
// against the Kubernetes API.
package service

import (
	"context"

	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/client/k8srestconfig"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/service/deployer"
	httpspec "github.com/giantswarm/draughtsman/service/http"
	slackspec "github.com/giantswarm/draughtsman/service/slack"
	"github.com/giantswarm/draughtsman/service/version"
)

// Config represents the configuration used to create a new service.
type Config struct {
	// Dependencies.
	FileSystem  afero.Fs
	HTTPClient  httpspec.Client
	Logger      micrologger.Logger
	SlackClient slackspec.Client

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper

	Description string
	GitCommit   string
	Name        string
	Source      string
}

// DefaultConfig provides a default configuration to create a new service by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		FileSystem:  afero.NewMemMapFs(),
		HTTPClient:  nil,
		Logger:      nil,
		SlackClient: nil,

		// Settings.
		Flag:  nil,
		Viper: nil,

		Description: "",
		GitCommit:   "",
		Name:        "",
		Source:      "",
	}
}

// New creates a new configured service object.
func New(config Config) (*Service, error) {
	// Settings.
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "viper must not be empty")
	}

	var err error

	var restConfig *rest.Config
	{
		c := k8srestconfig.Config{
			Logger: config.Logger,

			Address:   config.Viper.GetString(config.Flag.Service.Kubernetes.Address),
			InCluster: config.Viper.GetBool(config.Flag.Service.Kubernetes.InCluster),
			TLS: k8srestconfig.TLSClientConfig{
				CAFile:  config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CaFile),
				CrtFile: config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CrtFile),
				KeyFile: config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.KeyFile),
			},
		}

		restConfig, err = k8srestconfig.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	k8sClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var deployerService deployer.Deployer
	{
		deployerConfig := deployer.DefaultConfig()

		deployerConfig.FileSystem = config.FileSystem
		deployerConfig.HTTPClient = config.HTTPClient
		deployerConfig.KubernetesClient = k8sClient
		deployerConfig.Logger = config.Logger
		deployerConfig.SlackClient = config.SlackClient

		deployerConfig.Flag = config.Flag
		deployerConfig.Viper = config.Viper

		deployerConfig.Type = deployer.DeployerType(config.Viper.GetString(config.Flag.Service.Deployer.Type))

		deployerService, err = deployer.New(deployerConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var helmClient helmclient.Interface
	{
		c := helmclient.Config{
			K8sClient: k8sClient,
			Logger:    config.Logger,

			RestConfig:      restConfig,
			TillerNamespace: metav1.NamespaceSystem,
		}

		helmClient, err = helmclient.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var versionService *version.Service
	{
		versionConfig := version.DefaultConfig()

		versionConfig.Description = config.Description
		versionConfig.GitCommit = config.GitCommit
		versionConfig.Name = config.Name
		versionConfig.Source = config.Source

		versionService, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newService := &Service{
		// Dependencies.
		Deployer:   deployerService,
		HelmClient: helmClient,
		Version:    versionService,
	}

	return newService, nil
}

type Service struct {
	// Dependencies.
	Deployer   deployer.Deployer
	HelmClient helmclient.Interface
	Version    *version.Service
}

func (s *Service) Boot(ctx context.Context) {
	s.HelmClient.EnsureTillerInstalled(ctx)
	s.Deployer.Boot()
}
