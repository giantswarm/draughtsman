package service

import (
	"github.com/giantswarm/draughtsman/flag/service/deployer"
	"github.com/giantswarm/draughtsman/flag/service/kubernetes"
)

type Service struct {
	Deployer deployer.Deployer

	HTTPClientTimeout string

	Kubernetes kubernetes.Kubernetes

	SlackToken string
}
