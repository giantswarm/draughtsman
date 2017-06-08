package installer

import (
	"github.com/giantswarm/draughtsman/flag/service/deployer/installer/configurer"
	"github.com/giantswarm/draughtsman/flag/service/deployer/installer/helm"
)

type Installer struct {
	Helm       helm.Helm
	Configurer configurer.Configurer
	Type       string
}
