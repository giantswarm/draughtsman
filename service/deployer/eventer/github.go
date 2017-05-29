package eventer

import (
	micrologger "github.com/giantswarm/microkit/logger"
)

// GithubEventer is an Eventer that uses Github Deployment Events as a backend.
var GithubEventer EventerType = "GithubEventer"

// githubEventer is an implementer of the Eventer interface,
// that uses GitHub Deployment Events as a backend.
type githubEventer struct {
	// Dependencies.
	logger micrologger.Logger
}

func (e *githubEventer) NewDeploymentEvents() (<-chan DeploymentEvent, error) {
	e.logger.Log("debug", "starting polling for github deployment events")

	deploymentEventChannel := make(chan DeploymentEvent)

	return deploymentEventChannel, nil
}

func (e *githubEventer) SetPending(event DeploymentEvent) error {
	return nil
}

func (e *githubEventer) SetSuccess(event DeploymentEvent) error {
	return nil
}

func (e *githubEventer) SetFailed(event DeploymentEvent) error {
	return nil
}
