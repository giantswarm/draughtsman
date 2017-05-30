package github

import (
	"github.com/juju/errgo"
)

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var unexpectedStatusCode = errgo.New("unexpected status code")

// IsUnexpectedStatusCode asserts unexpectedStatusCode.
func IsUnexpectedStatusCode(err error) bool {
	return errgo.Cause(err) == unexpectedStatusCode
}
