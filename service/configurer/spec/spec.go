package spec

// ConfigurerType represents the type of Configurer to configure.
type ConfigurerType string

// Configurer represents a Service that provides a Helm configuration file.
type Configurer interface {
	// File returns a path to a file to use for Helm values.
	File() (string, error)
}
