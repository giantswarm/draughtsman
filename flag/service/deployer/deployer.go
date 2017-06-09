package deployer

import (
	"github.com/giantswarm/draughtsman/flag/service/deployer/eventer"
	"github.com/giantswarm/draughtsman/flag/service/deployer/installer"
	"github.com/giantswarm/draughtsman/flag/service/deployer/notifier"
)

type Deployer struct {
	Environment string
	Eventer     eventer.Eventer
	Installer   installer.Installer
	Notifier    notifier.Notifier
	Type        string
}
