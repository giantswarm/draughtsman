package flag

import (
	"github.com/giantswarm/microkit/flag"

	"github.com/giantswarm/draughtsman/flag/release"
	"github.com/giantswarm/draughtsman/flag/service"
)

type Flag struct {
	Release release.Release
	Service service.Service
}

func New() *Flag {
	f := &Flag{}
	flag.Init(f)
	return f
}
