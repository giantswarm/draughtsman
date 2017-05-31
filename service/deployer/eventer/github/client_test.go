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

// TestFilterDeploymentsByStatus tests the filterDeploymentsByStatus method.
func TestFilterDeploymentsByStatus(t *testing.T) {
	tests := []struct {
		deployments         []deployment
		expectedDeployments []deployment
	}{
		// Test that no deployments filters to no deployments.
		{
			deployments:         []deployment{},
			expectedDeployments: []deployment{},
		},

		// Test that a pending deployment is kept.
		{
			deployments: []deployment{
				deployment{
					Statuses: []deploymentStatus{
						deploymentStatus{State: pendingState},
					},
				},
			},
			expectedDeployments: []deployment{
				deployment{
					Statuses: []deploymentStatus{
						deploymentStatus{State: pendingState},
					},
				},
			},
		},

		// Test that a success only deployment is not kept.
		{
			deployments: []deployment{
				deployment{
					Statuses: []deploymentStatus{
						deploymentStatus{State: successState},
					},
				},
			},
			expectedDeployments: []deployment{},
		},

		// Test that a failure only deployment is not kept.
		{
			deployments: []deployment{
				deployment{
					Statuses: []deploymentStatus{
						deploymentStatus{State: failedState},
					},
				},
			},
			expectedDeployments: []deployment{},
		},

		// Test that a deployment that was pending, and is now successful, is not kept.
		{
			deployments: []deployment{
				deployment{
					Statuses: []deploymentStatus{
						deploymentStatus{State: pendingState},
						deploymentStatus{State: successState},
					},
				},
			},
			expectedDeployments: []deployment{},
		},

		// Test that a deployment that was pending, and is now failed, is not kept.
		{
			deployments: []deployment{
				deployment{
					Statuses: []deploymentStatus{
						deploymentStatus{State: pendingState},
						deploymentStatus{State: failedState},
					},
				},
			},
			expectedDeployments: []deployment{},
		},
	}

	for index, test := range tests {
		e := GithubEventer{}

		returnedDeployments := e.filterDeploymentsByStatus(test.deployments)

		if !reflect.DeepEqual(test.expectedDeployments, returnedDeployments) {
			t.Fatalf(
				"%v\nexpected: %#v\nreturned: %#v\n",
				index, test.expectedDeployments, returnedDeployments,
			)
		}
	}
}
