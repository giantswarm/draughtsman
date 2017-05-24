package flag

import (
	"github.com/giantswarm/microkit/flag"
)

type Flag struct {
}

func New() *Flag {
	f := &Flag{}
	flag.Init(f)
	return f
}
