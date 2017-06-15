package main

import (
	"net/http"
	"os"
	"time"

	"github.com/nlopes/slack"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/microkit/command"
	"github.com/giantswarm/microkit/logger"
	microserver "github.com/giantswarm/microkit/server"
	"github.com/giantswarm/operatorkit/client/k8s"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/server"
	"github.com/giantswarm/draughtsman/service"
	"github.com/giantswarm/draughtsman/service/deployer"
	"github.com/giantswarm/draughtsman/service/deployer/eventer/github"
	"github.com/giantswarm/draughtsman/service/deployer/installer/configurer/configmap"
	"github.com/giantswarm/draughtsman/service/deployer/installer/helm"
	slacknotifier "github.com/giantswarm/draughtsman/service/deployer/notifier/slack"
	slackspec "github.com/giantswarm/draughtsman/slack"
)

var (
	description string     = "draughtsman is an in-cluster agent that handles Helm based deployments."
	f           *flag.Flag = flag.New()
	gitCommit   string     = "n/a"
	name        string     = "draughtsman"
	source      string     = "https://github.com/giantswarm/draughtsman"
)

func main() {
	var err error

	var newLogger logger.Logger
	{
		loggerConfig := logger.DefaultConfig()
		loggerConfig.IOWriter = os.Stdout
		newLogger, err = logger.New(loggerConfig)
		if err != nil {
			panic(err)
		}
	}

	// We define a server factory to create the custom server once all command
	// line flags are parsed and all microservice configuration is storted out.
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

		var newKubernetesClient kubernetes.Interface
		{
			k8sConfig := k8s.Config{
				Logger: newLogger,

				Address:   v.GetString(f.Service.Kubernetes.Address),
				InCluster: v.GetBool(f.Service.Kubernetes.InCluster),
				TLS: k8s.TLSClientConfig{
					CAFile:  v.GetString(f.Service.Kubernetes.TLS.CaFile),
					CrtFile: v.GetString(f.Service.Kubernetes.TLS.CrtFile),
					KeyFile: v.GetString(f.Service.Kubernetes.TLS.KeyFile),
				},
			}
			newKubernetesClient, err = k8s.NewClient(k8sConfig)
			if err != nil {
				panic(err)
			}
		}

		var newSlackClient slackspec.Client
		{
			newSlackClient = slack.New(v.GetString(f.Service.Slack.Token))
		}

		// Create a new custom service which implements business logic.
		var newService *service.Service
		{
			serviceConfig := service.DefaultConfig()

			serviceConfig.HTTPClient = newHttpClient
			serviceConfig.KubernetesClient = newKubernetesClient
			serviceConfig.Logger = newLogger
			serviceConfig.SlackClient = newSlackClient

			serviceConfig.Flag = f
			serviceConfig.Viper = v

			serviceConfig.Description = description
			serviceConfig.GitCommit = gitCommit
			serviceConfig.Name = name
			serviceConfig.Source = source

			newService, err = service.New(serviceConfig)
			if err != nil {
				panic(err)
			}
		}

		// Create a new custom server which bundles our endpoints.
		var newServer microserver.Server
		{
			serverConfig := server.DefaultConfig()

			serverConfig.MicroServerConfig.Logger = newLogger
			serverConfig.MicroServerConfig.ServiceName = name
			serverConfig.MicroServerConfig.Viper = v
			serverConfig.Service = newService

			newServer, err = server.New(serverConfig)
			if err != nil {
				panic(err)
			}
			go newService.Boot()
		}

		return newServer
	}

	// Create a new microkit command which manages our custom microservice.
	var newCommand command.Command
	{
		commandConfig := command.DefaultConfig()

		commandConfig.Logger = newLogger
		commandConfig.ServerFactory = newServerFactory

		commandConfig.Description = description
		commandConfig.GitCommit = gitCommit
		commandConfig.Name = name
		commandConfig.Source = source

		newCommand, err = command.New(commandConfig)
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
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Configurer.Type, string(configmap.ConfigmapConfigurerType), "Which configurer to use for configuration management.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Notifier.Type, string(slacknotifier.SlackNotifierType), "Which notifier to use for notification management.")

	// Client configuration.
	daemonCommand.PersistentFlags().Duration(f.Service.HTTPClient.Timeout, 10*time.Second, "Timeout for HTTP requests.")

	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.Address, "", "Address used to connect to Kubernetes. When empty in-cluster config is created.")
	daemonCommand.PersistentFlags().Bool(f.Service.Kubernetes.InCluster, true, "Whether to use the in-cluster config to authenticate with Kubernetes.")
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

	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Configurer.Configmap.Key, "values", "Key in configmap holding values data.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Configurer.Configmap.Name, "draughtsman-values", "Name of configmap holding values data.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Configurer.Configmap.Namespace, "draughtsman", "Namespace of configmap holding values data.")

	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Configurer.File.Path, "", "Path to values file.")

	daemonCommand.PersistentFlags().String(f.Service.Deployer.Notifier.Slack.Channel, "", "Channel to post Slack notifications to.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Notifier.Slack.Emoji, ":older_man:", "Emoji to use for Slack notifications.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Notifier.Slack.Username, "draughtsman", "Username to post Slack notifications with.")

	newCommand.CobraCommand().Execute()
}
