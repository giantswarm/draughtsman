package file

import (
	"os"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/draughtsman/service/deployer/installer/configurer/spec"
)

// FileConfigurerType is a Configurer that uses a normal file.
var FileConfigurerType spec.ConfigurerType = "FileConfigurer"

// Config represents the configuration used to create a File Configurer.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Path string
}

// DefaultConfig provides a default configuration to create a new File
// Configurer by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,
	}
}

// New creates a new configured File Configurer.
func New(config Config) (*FileConfigurer, error) {
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	if config.Path == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "path must not be empty")
	}

	if _, err := os.Stat(config.Path); os.IsNotExist(err) {
		return nil, microerror.MaskAnyf(invalidConfigError, "path does not exist")
	}

	configurer := &FileConfigurer{
		// Dependencies.
		logger: config.Logger,

		// Settings.
		path: config.Path,
	}

	return configurer, nil
}

// FileConfigurer is an implementation of the Configurer interface,
// that uses a plain file to hold configuration.
type FileConfigurer struct {
	// Dependencies.
	logger micrologger.Logger

	// Settings.
	path string
}

func (c *FileConfigurer) File() (string, error) {
	return c.path, nil
}
