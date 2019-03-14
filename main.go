package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/microkit/command"
	microserver "github.com/giantswarm/microkit/server"
	"github.com/giantswarm/micrologger"
	"github.com/nlopes/slack"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/server"
	"github.com/giantswarm/draughtsman/service"
	"github.com/giantswarm/draughtsman/service/configurer/configmap"
	"github.com/giantswarm/draughtsman/service/configurer/secret"
	"github.com/giantswarm/draughtsman/service/deployer"
	"github.com/giantswarm/draughtsman/service/eventer/github"
	"github.com/giantswarm/draughtsman/service/installer/helm"
	slacknotifier "github.com/giantswarm/draughtsman/service/notifier/slack"
	slackspec "github.com/giantswarm/draughtsman/service/slack"
)

var (
	description string     = "draughtsman is an in-cluster agent that handles Helm based deployments."
	f           *flag.Flag = flag.New()
	gitCommit   string     = "n/a"
	name        string     = "draughtsman"
	source      string     = "https://github.com/giantswarm/draughtsman"
)

func main() {
	err := mainError()
	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}
}

func mainError() error {
	var err error

	ctx := context.Background()
	newLogger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		return microerror.Mask(err)
	}

	// We define a server factory to create the custom server once all command
	// line flags are parsed and all microservice configuration is sorted out.
	newServerFactory := func(v *viper.Viper) microserver.Server {
		var newHttpClient *http.Client
		{
			httpClientTimeout := v.GetDuration(f.Service.HTTPClient.Timeout)
			if httpClientTimeout.Seconds() == 0 {
				panic("http client timeout must be greater than zero")
			}

			newHttpClient = &http.Client{
				Timeout: httpClientTimeout,
			}
		}

		var newSlackClient slackspec.Client
		{
			newSlackClient = slack.New(v.GetString(f.Service.Slack.Token))
		}

		// Create a new custom service which implements business logic.
		var newService *service.Service
		{
			c := service.Config{
				Flag:        f,
				FileSystem:  afero.NewOsFs(),
				HTTPClient:  newHttpClient,
				Logger:      newLogger,
				SlackClient: newSlackClient,
				Viper:       v,

				Description: description,
				GitCommit:   gitCommit,
				ProjectName: name,
				Source:      source,
			}

			newService, err = service.New(c)
			if err != nil {
				panic(fmt.Sprintf("%#v", microerror.Mask(err)))
			}

			go newService.Boot(ctx)
		}

		// New custom server that bundles microkit endpoints.
		var newServer microserver.Server
		{
			c := server.Config{
				Logger:      newLogger,
				Service:     newService,
				Viper:       v,
				ProjectName: name,
			}

			newServer, err = server.New(c)
			if err != nil {
				panic(fmt.Sprintf("%#v\n", microerror.Maskf(err, "server.New")))
			}
		}

		return newServer
	}

	// Create a new microkit command which manages our custom microservice.
	var newCommand command.Command
	{
		c := command.Config{
			Logger:        newLogger,
			ServerFactory: newServerFactory,

			Description: description,
			GitCommit:   gitCommit,
			Name:        name,
			Source:      source,
		}

		newCommand, err = command.New(c)
		if err != nil {
			panic(err)
		}
	}

	daemonCommand := newCommand.DaemonCommand().CobraCommand()

	daemonCommand.PersistentFlags().String(f.Service.Deployer.Environment, "", "Environment name that draughtsman is running in.")

	// Component type selection.
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Type, string(deployer.StandardDeployer), "Which deployer to use for deployment management.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Eventer.Type, string(github.GithubEventerType), "Which eventer to use for event management.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Type, string(helm.HelmInstallerType), "Which installer to use for installation management.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Configurer.Types, string(configmap.ConfigurerType)+","+string(secret.ConfigurerType), "Comma separated list of configurers to use for configuration management.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Notifier.Type, string(slacknotifier.SlackNotifierType), "Which notifier to use for notification management.")

	// Client configuration.
	daemonCommand.PersistentFlags().Duration(f.Service.HTTPClient.Timeout, 10*time.Second, "Timeout for HTTP requests.")

	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.Address, "", "Address used to connect to Kubernetes. When empty in-cluster config is created.")
	daemonCommand.PersistentFlags().Bool(f.Service.Kubernetes.InCluster, true, "Whether to use the in-cluster config to authenticate with Kubernetes.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.KubeConfig, "", "KubeConfig used to connect to Kubernetes. When empty other settings are used.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.TLS.CaFile, "", "Certificate authority file path to use to authenticate with Kubernetes.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.TLS.CrtFile, "", "Certificate file path to use to authenticate with Kubernetes.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.TLS.KeyFile, "", "Key file path to use to authenticate with Kubernetes.")

	daemonCommand.PersistentFlags().String(f.Service.Slack.Token, "", "Token to post Slack notifications with.")

	// Service configuration.
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Eventer.GitHub.OAuthToken, "", "OAuth token for authenticating against GitHub. Needs 'repo_deployment' scope.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Eventer.GitHub.Organisation, "", "Organisation under which to check for deployments.")
	daemonCommand.PersistentFlags().Duration(f.Service.Deployer.Eventer.GitHub.PollInterval, 1*time.Minute, "Interval to poll for new deployments.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Eventer.GitHub.ProjectList, "", "Comma seperated list of GitHub projects to check for deployments.")

	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Helm.HelmBinaryPath, "/bin/helm", "Path to Helm binary. Needs CNR registry plugin installed.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Helm.Organisation, "", "Organisation of Helm CNR registry.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Helm.Password, "", "Password for Helm CNR registry.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Helm.Registry, "quay.io", "URL for Helm CNR registry.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Helm.Username, "", "Username for Helm CNR registry.")

	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Configurer.ConfigMap.Key, "values", "Key in configmap holding values data.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Configurer.ConfigMap.Name, "draughtsman-values-configmap", "Name of configmap holding values data.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Configurer.ConfigMap.Namespace, "draughtsman", "Namespace of configmap holding values data.")

	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Configurer.File.Path, "", "Path to values file.")

	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Configurer.Secret.Key, "values", "Key in secret holding values data.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Configurer.Secret.Name, "draughtsman-values-secret", "Name of secret holding values data.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Configurer.Secret.Namespace, "draughtsman", "Namespace of secret holding values data.")

	daemonCommand.PersistentFlags().String(f.Service.Deployer.Notifier.Slack.Channel, "", "Channel to post Slack notifications to.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Notifier.Slack.Emoji, ":older_man:", "Emoji to use for Slack notifications.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Notifier.Slack.Username, "draughtsman", "Username to post Slack notifications with.")

	newCommand.CobraCommand().Execute()

	return nil
}
