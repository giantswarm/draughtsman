package configmap

import (
	"github.com/juju/errgo"
)

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var keyMissingError = errgo.New("key missing")

// IsKeyMissing asserts keyMissingError
func IsKeyMissing(err error) bool {
	return errgo.Cause(err) == keyMissingError
}
