package k8sclient

import "testing"

func Test_Interface(t *testing.T) {
	var _ Interface = &Clients{}
}
