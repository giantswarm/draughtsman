package github

import (
	"reflect"
	"testing"
)

// TestFilterDeploymentsByEnvironment tests the filterDeployments method.
func TestFilterDeploymentsByEnvironment(t *testing.T) {
	tests := []struct {
		deployments         []deployment
		environment         string
		expectedDeployments []deployment
	}{
		// Test that no deployments filters to no deployments.
		{
			deployments:         []deployment{},
			environment:         "covfefe",
			expectedDeployments: []deployment{},
		},

		// Test that a deployment for the installation is kept.
		{
			deployments: []deployment{
				deployment{Environment: "production"},
			},
			environment: "production",
			expectedDeployments: []deployment{
				deployment{Environment: "production"},
			},
		},

		// Test that only this environment's deployments are kept.
		{
			deployments: []deployment{
				deployment{Environment: "development"},
				deployment{Environment: "production"},
			},
			environment: "development",
			expectedDeployments: []deployment{
				deployment{Environment: "development"},
			},
		},
	}

	for index, test := range tests {
		e := GithubEventer{
			environment: test.environment,
		}

		returnedDeployments := e.filterDeploymentsByEnvironment(test.deployments)

		if !reflect.DeepEqual(test.expectedDeployments, returnedDeployments) {
			t.Fatalf(
				"%v\nexpected: %#v\nreturned: %#v\n",
				index, test.expectedDeployments, returnedDeployments,
			)
		}
	}
}
