package deployer

import (
	"github.com/giantswarm/draughtsman/flag/service/deployer/eventer"
	"github.com/giantswarm/draughtsman/flag/service/deployer/installer"
)

type Deployer struct {
	Eventer   eventer.Eventer
	Installer installer.Installer
}
