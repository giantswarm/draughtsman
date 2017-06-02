package main

import (
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/giantswarm/microkit/command"
	"github.com/giantswarm/microkit/logger"
	microserver "github.com/giantswarm/microkit/server"

	"github.com/giantswarm/draughtsman/flag"
	"github.com/giantswarm/draughtsman/server"
	"github.com/giantswarm/draughtsman/service"
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

	// Create a new logger which is used by all packages.
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
		// Create a new custom service which implements business logic.
		var newService *service.Service
		{
			serviceConfig := service.DefaultConfig()

			serviceConfig.Logger = newLogger

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

	daemonCommand.PersistentFlags().String(f.Service.Deployer.Eventer.GitHub.Environment, "", "Environment name that draughtsman is running in.")
	daemonCommand.PersistentFlags().Duration(f.Service.Deployer.Eventer.GitHub.HTTPClientTimeout, 10*time.Second, "Timeout for requests to GitHub.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Eventer.GitHub.OAuthToken, "", "OAuth token for authenticating against GitHub. Needs 'repo_deployment' scope.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Eventer.GitHub.Organisation, "", "Organisation under which to check for deployments.")
	daemonCommand.PersistentFlags().Duration(f.Service.Deployer.Eventer.GitHub.PollInterval, 1*time.Minute, "Interval to poll for new deployments.")
	daemonCommand.PersistentFlags().StringSlice(f.Service.Deployer.Eventer.GitHub.ProjectList, []string{}, "List of GitHub projects to check for deployments.")

	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Helm.HelmImage, "quay.io/giantswarm/docker-helm", "Docker image to run Helm with. Image needs CNR registry plugin installed.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Helm.HelmImageTag, "006b0db51ec484be8b1bd49990804784a9737ece", "Tag of the Helm image.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Helm.Password, "", "Password for Helm CNR registry.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Helm.Registry, "quay.io", "URL for Helm CNR registry.")
	daemonCommand.PersistentFlags().String(f.Service.Deployer.Installer.Helm.Username, "", "Username for Helm CNR registry.")

	newCommand.CobraCommand().Execute()
}
