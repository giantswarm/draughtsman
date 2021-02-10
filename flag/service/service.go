package service

import (
	"github.com/giantswarm/operatorkit/flag/service/kubernetes"

	"github.com/giantswarm/draughtsman/flag/service/deployer"
	"github.com/giantswarm/draughtsman/flag/service/helmmigration"
	"github.com/giantswarm/draughtsman/flag/service/httpclient"
	"github.com/giantswarm/draughtsman/flag/service/slack"
)

type Service struct {
	Deployer      deployer.Deployer
	HelmMigration helmmigration.HelmMigration

	HTTPClient httpclient.HTTPClient
	Kubernetes kubernetes.Kubernetes
	Slack      slack.Slack
}
