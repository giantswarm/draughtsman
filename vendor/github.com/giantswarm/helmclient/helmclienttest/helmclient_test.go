package helmclienttest

import (
	"testing"
)

func Test_New(t *testing.T) {
	// Test that New doesn't panic and helmclient.Interface is implemented.
	New(Config{})
}
