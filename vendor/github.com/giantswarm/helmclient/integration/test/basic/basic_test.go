// +build k8srequired

package basic

import (
	"context"
	"os"
	"testing"

	"github.com/giantswarm/helmclient/integration/charttarball"
	"k8s.io/helm/pkg/helm"
)

func TestBasic(t *testing.T) {
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

	// We need to pass the ValueOverrides option to make the install process
	// use the default values and prevent errors on nested values.
	//
	//     {
	//      rpc error: code = Unknown desc = render error in "cnr-server-chart/templates/deployment.yaml":
	//      template: cnr-server-chart/templates/deployment.yaml:20:26:
	//      executing "cnr-server-chart/templates/deployment.yaml" at <.Values.image.reposi...>: can't evaluate field repository in type interface {}
	//     }
	//
	err = config.HelmClient.InstallReleaseFromTarball(ctx, tarballPath, "default", helm.ReleaseName(releaseName), helm.ValueOverrides([]byte("{}")))
	if err != nil {
		t.Fatalf("could not install chart %v", err)
	}

	releaseContent, err := config.HelmClient.GetReleaseContent(ctx, releaseName)
	if err != nil {
		t.Fatalf("could not get release content %v", err)
	}

	expectedName := releaseName
	actualName := releaseContent.Name
	if expectedName != actualName {
		t.Fatalf("bad release name, want %q, got %q", expectedName, actualName)
	}

	expectedStatus := "DEPLOYED"
	actualStatus := releaseContent.Status
	if expectedStatus != actualStatus {
		t.Fatalf("bad release status, want %q, got %q", expectedStatus, actualStatus)
	}

	err = config.HelmClient.RunReleaseTest(ctx, releaseName)
	if err != nil {
		t.Fatalf("error running tests, want nil got %v", err)
	}

	// Test should fail on the 2nd attempt because the test pod already exists.
	err = config.HelmClient.RunReleaseTest(ctx, releaseName)
	if err == nil {
		t.Fatalf("error running tests, want error got nil")
	}

	err = config.HelmClient.DeleteRelease(ctx, releaseName)
	if err != nil {
		t.Fatalf("could not delete release %v", err)
	}

	releaseContent, err = config.HelmClient.GetReleaseContent(ctx, releaseName)
	if err != nil {
		t.Fatalf("could not get release content %v", err)
	}
	expectedStatus = "DELETED"
	actualStatus = releaseContent.Status
	if expectedStatus != actualStatus {
		t.Fatalf("bad release status, want %q, got %q", expectedStatus, actualStatus)
	}
}
