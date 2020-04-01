package eventer

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

var unexpectedStatusCode = &microerror.Error{
	Kind: "unexpectedStatusCodeError",
}

// IsUnexpectedStatusCode asserts unexpectedStatusCode.
func IsUnexpectedStatusCode(err error) bool {
	return microerror.Cause(err) == unexpectedStatusCode
}
