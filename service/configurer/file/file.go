package file

import (
	"os"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"

	"github.com/giantswarm/draughtsman/service/configurer/spec"
)

// ConfigurerType is the kind of a Configurer that uses a normal file.
var ConfigurerType spec.ConfigurerType = "FileConfigurer"

// Config represents the configuration used to create a File Configurer.
type Config struct {
	// Dependencies.
	FileSystem afero.Fs
	Logger     micrologger.Logger

	// Settings.
	Path string
}

// DefaultConfig provides a default configuration to create a new File
// Configurer by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		FileSystem: afero.NewMemMapFs(),
		Logger:     nil,

		// Settings.
		Path: "",
	}
}

// New creates a new configured File Configurer.
func New(config Config) (*FileConfigurer, error) {
	// Dependencies.
	if config.FileSystem == nil {
		return nil, microerror.Maskf(invalidConfigError, "file system must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "logger must not be empty")
	}

	// Settings.
	if config.Path == "" {
		return nil, microerror.Maskf(invalidConfigError, "path must not be empty")
	}

	_, err := os.Stat(config.Path)
	if os.IsNotExist(err) {
		return nil, microerror.Maskf(invalidConfigError, "path does not exist")
	}

	configurer := &FileConfigurer{
		// Dependencies.
		fileSystem: config.FileSystem,
		logger:     config.Logger,

		// Settings.
		path: config.Path,
	}

	return configurer, nil
}

// FileConfigurer is an implementation of the Configurer interface,
// that uses a plain file to hold configuration.
type FileConfigurer struct {
	// Dependencies.
	fileSystem afero.Fs
	logger     micrologger.Logger

	// Settings.
	path string
}

func (c *FileConfigurer) Type() spec.ConfigurerType {
	return ConfigurerType
}

func (c *FileConfigurer) Values() (string, error) {
	b, err := afero.ReadFile(c.fileSystem, c.path)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return string(b), nil
}
