package eventer

import (
	"github.com/giantswarm/draughtsman/flag/service/deployer/eventer/github"
)

type Eventer struct {
	GitHub github.GitHub
	Type   string
}
