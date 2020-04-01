package slack

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

var slackError = &microerror.Error{
	Kind: "slackError",
}

// IsSlackError asserts slackError.
func IsSlackError(err error) bool {
	return microerror.Cause(err) == slackError
}
