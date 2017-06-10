package service

import (
	"github.com/giantswarm/draughtsman/flag/service/deployer"
)

type Service struct {
	Deployer          deployer.Deployer
	HTTPClientTimeout string
}
