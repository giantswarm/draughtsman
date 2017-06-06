package helm

import (
	"github.com/juju/errgo"
)

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var helmError = errgo.New("helm error")

// IsHelm asserts helmError.
func IsHelm(err error) bool {
	return errgo.Cause(err) == helmError
}
