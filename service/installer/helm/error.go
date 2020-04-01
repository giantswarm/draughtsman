package helm

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var helmError = &microerror.Error{
	Kind: "helmError",
}

// IsHelm asserts helmError.
func IsHelm(err error) bool {
	return microerror.Cause(err) == helmError
}
