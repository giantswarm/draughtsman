package helm

import (
	"testing"
)

// TestVersionedChartName tests the versionedChartName method.
func TestVersionedChartName(t *testing.T) {
	tests := []struct {
		registry          string
		organisation      string
		project           string
		sha               string
		expectedChartName string
	}{
		{
			registry:          "quay.io",
			organisation:      "giantswarm",
			project:           "api",
			sha:               "12345",
			expectedChartName: "quay.io/giantswarm/api-chart@1.0.0-12345",
		},
	}

	for index, test := range tests {
		i := HelmInstaller{
			registry:     test.registry,
			organisation: test.organisation,
		}

		returnedChartName := i.versionedChartName(test.project, test.sha)

		if returnedChartName != test.expectedChartName {
			t.Fatalf(
				"%v\nexpected: %#v\nreturned: %#v\n",
				index, test.expectedChartName, returnedChartName,
			)
		}
	}
}

// TestTarballName tests the tarballName method.
func TestTarballName(t *testing.T) {
	tests := []struct {
		registry            string
		organisation        string
		project             string
		sha                 string
		expectedTarballName string
	}{
		{
			registry:            "quay.io",
			organisation:        "giantswarm",
			project:             "api",
			sha:                 "12345",
			expectedTarballName: "giantswarm_api-chart_1.0.0-12345.tar.gz",
		},
	}

	for index, test := range tests {
		i := HelmInstaller{
			registry:     test.registry,
			organisation: test.organisation,
		}

		returnedTarballName := i.tarballName(test.project, test.sha)

		if returnedTarballName != test.expectedTarballName {
			t.Fatalf(
				"%v\nexpected: %#v\nreturned: %#v\n",
				index, test.expectedTarballName, returnedTarballName,
			)
		}
	}
}
