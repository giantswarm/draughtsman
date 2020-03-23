package github

import (
	"github.com/giantswarm/microerror"
)

var executionFailedError = &microerror.Error{
	Kind: "executionFailedError",
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var missingHeaderError = &microerror.Error{
	Kind: "missingHeaderError",
}

// IsMissingHeaderError asserts missingHeaderError.
func IsMissingHeaderError(err error) bool {
	return microerror.Cause(err) == missingHeaderError
}

var unexpectedStatusCode = &microerror.Error{
	Kind: "unexpectedStatusCode",
}

// IsUnexpectedStatusCode asserts unexpectedStatusCode.
func IsUnexpectedStatusCode(err error) bool {
	return microerror.Cause(err) == unexpectedStatusCode
}
