package spec

// ConfigurerType represents the type of Configurer to configure.
type ConfigurerType string

// Configurer represents a Service that provides a Helm configuration file.
type Configurer interface {
	// Type returns the configurer type of the current implementation.
	Type() ConfigurerType
	// Values returns content of a file to use for Helm values. The caller is
	// responsible for persisting and cleaning up eventual files on the file
	// system.
	Values() (string, error)
}
