package deployer

import (
	"github.com/giantswarm/draughtsman/flag/service/deployer/eventer"
)

type Deployer struct {
	Eventer eventer.Eventer
}
