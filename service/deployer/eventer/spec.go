package eventer

import (
	"time"
)

// DeploymentEvent represents a request for a package to be deployed.
type DeploymentEvent struct {
	// Name is the name of the project to deploy, e.g: aws-operator.
	Name string
}

// Eventer represents a Service that checks for deployment events.
type Eventer interface {
	// NewDeploymentEvents returns a list of DeploymentEvents that have been
	// created since the previous check, and a minimum time for the controller
	// to wait before checking again.
	// In case of error, the error will be non-nil.
	NewDeploymentEvents() ([]DeploymentEvent, time.Duration, error)

	// SetPending updates the DeploymentEvent remote state to a pending state.
	SetPending(DeploymentEvent) error
	// SetSuccess updates the DeploymentEvent remote state to a success state.
	SetSuccess(DeploymentEvent) error
	// SetFailed updates the DeploymentEvent remote state to a failed state.
	SetFailed(DeploymentEvent) error
}
