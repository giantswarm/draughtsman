// +build k8srequired

package basic

import (
	"context"
	"os"
	"testing"

	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/helmclient/integration/charttarball"
)

func TestDeleteRelease_IsReleaseNotFound(t *testing.T) {
	ctx := context.Background()
	var err error

	const releaseName = "test"

	tarballPath, err := charttarball.Create("test-chart")
	if err != nil {
		t.Fatalf("could not create chart archive %#v", err)
	}
	defer os.Remove(tarballPath)

	err = config.HelmClient.EnsureTillerInstalled(ctx)
	if err != nil {
		t.Fatalf("could not install Tiller %#v", err)
	}

	err = config.HelmClient.DeleteRelease(ctx, releaseName)
	if helmclient.IsReleaseNotFound(err) {
		// This is error we want.
	} else if err != nil {
		t.Fatalf("unexpected error while deleting release %#v", err)
	} else {
		t.Fatalf("expected error for already deleted release, got %#v", err)
	}
}
