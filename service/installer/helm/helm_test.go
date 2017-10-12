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

// TestChartName tests the chartName method.
func TestChartName(t *testing.T) {
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
			expectedChartName: "giantswarm_api-chart_1.0.0-12345/api-chart",
		},
	}

	for index, test := range tests {
		i := HelmInstaller{
			registry:     test.registry,
			organisation: test.organisation,
		}

		returnedChartName := i.chartName(test.project, test.sha)

		if returnedChartName != test.expectedChartName {
			t.Fatalf(
				"%v\nexpected: %#v\nreturned: %#v\n",
				index, test.expectedChartName, returnedChartName,
			)
		}
	}
}
