package configurer

import (
	"github.com/giantswarm/draughtsman/flag/service/deployer/installer/configurer/configmap"
	"github.com/giantswarm/draughtsman/flag/service/deployer/installer/configurer/file"
)

type Configurer struct {
	Configmap configmap.Configmap
	File      file.File
	Type      string
}
