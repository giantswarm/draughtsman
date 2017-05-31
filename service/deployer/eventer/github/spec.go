package github

import (
	"github.com/giantswarm/draughtsman/service/deployer/eventer/spec"
)

// deployment represents a GitHub API Deployment.
// See: https://developer.github.com/v3/repos/deployments/#create-a-deployment
type deployment struct {
	// Environment is the environment field of the GitHub deployment.
	Environment string `json:"environment"`

	// ID is the ID field of the GitHub deployment.
	ID int `json:"id"`

	// Sha is the SHA hash of the commit the deployment references.
	Sha string `json:"sha"`

	// Statuses is the deployment statuses of this deployment.
	Statuses []deploymentStatus
}

// Deploymespec.ntEvent returns the deployment as a DeploymentEvent.
func (d deployment) DeploymentEvent(project string) spec.DeploymentEvent {
	return spec.DeploymentEvent{
		ID:   d.ID,
		Name: project,
		Sha:  d.Sha,
	}
}

// deploymentStatus represents a GitHub API Deployment Status.
// See: https://developer.github.com/v3/repos/deployments/#create-a-deployment-status
type deploymentStatus struct {
	// State is the state of the deployment status.
	State deploymentStatusState `json:"state"`
}

// deploymentStatusState represents possible Deployment Status states.
type deploymentStatusState string

var (
	// pendingState is the state for pending Deployment Status states.
	pendingState deploymentStatusState = "pending"
	// successState is the state for successful Deployment Status states.
	successState deploymentStatusState = "success"
	// failedState is the state for failed Deployment Status states.
	failedState deploymentStatusState = "failed"
)
