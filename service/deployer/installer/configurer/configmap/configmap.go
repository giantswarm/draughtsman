package configmap

import (
	"io/ioutil"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	"github.com/giantswarm/operatorkit/client/k8s"

	"github.com/giantswarm/draughtsman/service/deployer/installer/configurer/spec"
)

// ConfigmapConfigurerType is a Configurer that is backed by a Kubernetes
// Configmap.
var ConfigmapConfigurerType spec.ConfigurerType = "ConfigmapConfigurer"

// Config represents the configuration used to create a Configmap Configurer.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Address     string
	CAFilePath  string
	CrtFilePath string
	InCluster   bool
	// Key is the key to reference the values data in the configmap.
	Key         string
	KeyFilePath string
	Name        string
	Namespace   string
}

// DefaultConfig provides a default configuration to create a new Configmap
// Configurer by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,
	}
}

// New creates a new configured Configmap Configurer.
func New(config Config) (*ConfigmapConfigurer, error) {
	if !config.InCluster {
		if config.Address == "" {
			return nil, microerror.MaskAnyf(invalidConfigError, "address must not be empty")
		}
	}

	if config.Key == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "key must not be empty")
	}
	if config.Name == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "name must not be empty")
	}
	if config.Namespace == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "namespace must not be empty")
	}

	k8sConfig := k8s.Config{
		Logger: config.Logger,

		Address:   config.Address,
		InCluster: config.InCluster,
		TLS: k8s.TLSClientConfig{
			CAFile:  config.CAFilePath,
			CrtFile: config.CrtFilePath,
			KeyFile: config.KeyFilePath,
		},
	}
	client, err := k8s.NewClient(k8sConfig)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	config.Logger.Log("debug", "checking connection to Kubernetes")
	if _, err := client.CoreV1().Namespaces().Get("default", v1.GetOptions{}); err != nil {
		return nil, microerror.MaskAny(err)
	}

	// Create a temporary file to use for holding the values file for Helm to read.
	tempFile, err := ioutil.TempFile("", "draughtsman")
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	configurer := &ConfigmapConfigurer{
		// Dependencies.
		client: client,
		logger: config.Logger,

		// Settings.
		key:       config.Key,
		name:      config.Name,
		namespace: config.Namespace,
		tempFile:  tempFile,
	}

	return configurer, nil
}

// ConfigmapConfigurer is an implementation of the Configurer interface,
// that uses a Kubernetes Configmap to hold configuration.
type ConfigmapConfigurer struct {
	// Dependencies.
	client kubernetes.Interface
	logger micrologger.Logger

	// Settings.
	key       string
	name      string
	namespace string
	tempFile  *os.File
}

func (c *ConfigmapConfigurer) File() (string, error) {
	defer updateConfigmapMetrics(time.Now())

	c.logger.Log("debug", "fetching configuration from configmap", "name", c.name, "namespace", c.namespace)

	cm, err := c.client.CoreV1().ConfigMaps(c.namespace).Get(c.name, v1.GetOptions{})
	if err != nil {
		return "", microerror.MaskAny(err)
	}

	valuesData, ok := cm.Data[c.key]
	if !ok {
		return "", microerror.MaskAnyf(keyMissingError, "key '%v' not found in configmap", c.key)
	}

	c.logger.Log("debug", "writing configuration to temp file", "path", c.tempFile.Name())

	if _, err := c.tempFile.WriteString(valuesData); err != nil {
		return "", microerror.MaskAny(err)
	}

	return c.tempFile.Name(), nil
}
