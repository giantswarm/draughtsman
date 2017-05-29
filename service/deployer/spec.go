package deployer

// Deployer is a service that handles deployments.
// While the deployment control loop logic is unlikely to change drastically,
// the Deployer interface exists to allow for it in the future.
// e.g: installations that only allow for deployments during certain
// maintenance windows.
type Deployer interface {
	Boot()
}
