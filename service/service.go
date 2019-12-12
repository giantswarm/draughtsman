// Package service implements business logic to create Kubernetes resources
// against the Kubernetes API.
package service

import (
	"context"
	"fmt"

	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/k8sclient/k8srestconfig"
	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/service/deployer"
	httpspec "github.com/giantswarm/draughtsman/service/http"
	slackspec "github.com/giantswarm/draughtsman/service/slack"
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
	ProjectName string
	Source      string
	Version     string
}

type Service struct {
	// Dependencies.
	Deployer   deployer.Deployer
	HelmClient helmclient.Interface
	Version    *version.Service
}

// New creates a new configured service object.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.FileSystem == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.FileSystem must not be empty", config)
	}
	if config.HTTPClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.HTTPClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.SlackClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.SlackClient must not be empty", config)
	}

	// Settings.
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flag must not be empty", config)
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "viper must not be empty")
	}

	var err error

	var restConfig *rest.Config
	{
		c := k8srestconfig.Config{
			Logger: config.Logger,

			Address:    config.Viper.GetString(config.Flag.Service.Kubernetes.Address),
			InCluster:  config.Viper.GetBool(config.Flag.Service.Kubernetes.InCluster),
			KubeConfig: config.Viper.GetString(config.Flag.Service.Kubernetes.KubeConfig),
			TLS: k8srestconfig.ConfigTLS{
				CAFile:  config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CAFile),
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
			// TODO: Remove once app CRs are used in all control plane
			// clusters. As chart-operator will then take care of upgrading
			// Tiller.
			//
			//	https://github.com/giantswarm/giantswarm/issues/8068
			//
			TillerUpgradeEnabled: true,
		}

		helmClient, err = helmclient.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var versionService *version.Service
	{
		versionConfig := version.Config{
			Description: config.Description,
			GitCommit:   config.GitCommit,
			Name:        config.ProjectName,
			Source:      config.Source,
			Version:     config.Version,
		}

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

func (s *Service) Boot(ctx context.Context) error {
	// TODO: Improve error handling in Boot method.
	//
	//	See https://github.com/giantswarm/giantswarm/issues/4965
	//
	err := s.HelmClient.EnsureTillerInstalled(ctx)
	if err != nil {
		panic(fmt.Sprintf("%#v\n", microerror.Maskf(err, "service.Boot")))
	}

	s.Deployer.Boot()

	return nil
}
