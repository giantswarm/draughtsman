package configurer

import (
	"github.com/giantswarm/draughtsman/flag/service/deployer/installer/configurer/configmap"
	"github.com/giantswarm/draughtsman/flag/service/deployer/installer/configurer/file"
	"github.com/giantswarm/draughtsman/flag/service/deployer/installer/configurer/secret"
)

type Configurer struct {
	Configmap configmap.Configmap
	File      file.File
	Secret    secret.Secret
	Type      string
}
