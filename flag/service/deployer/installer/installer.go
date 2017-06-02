package installer

import (
	"github.com/giantswarm/draughtsman/flag/service/deployer/installer/helm"
)

type Installer struct {
	Helm helm.Helm
}
