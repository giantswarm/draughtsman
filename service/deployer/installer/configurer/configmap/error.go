package configmap

import (
	"github.com/juju/errgo"
)

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var kubernetesError = errgo.New("kubernetes")

// IsKubernetes asserts kubernetesError.
func IsKubernetes(err error) bool {
	return errgo.Cause(err) == kubernetesError
}
