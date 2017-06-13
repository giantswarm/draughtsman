package service

import (
	"github.com/giantswarm/draughtsman/flag/service/deployer"
	"github.com/giantswarm/draughtsman/flag/service/kubernetes"
	"github.com/giantswarm/draughtsman/flag/service/slack"
)

type Service struct {
	Deployer deployer.Deployer

	HTTPClientTimeout string

	Kubernetes kubernetes.Kubernetes
	Slack      slack.Slack
}
