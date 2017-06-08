package configurer

import (
	"github.com/giantswarm/draughtsman/flag/service/deployer/installer/configurer/file"
)

type Configurer struct {
	File file.File
	Type string
}
