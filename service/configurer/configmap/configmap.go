package configmap

import (
	"io/ioutil"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/draughtsman/service/configurer/spec"
)

// ConfigMapConfigurerType is a Configurer that is backed by a Kubernetes
// ConfigMap.
var ConfigMapConfigurerType spec.ConfigurerType = "ConfigMapConfigurer"

// Config represents the configuration used to create a ConfigMap Configurer.
type Config struct {
	// Dependencies.
	KubernetesClient kubernetes.Interface
	Logger           micrologger.Logger

	// Settings.
	// Key is the key to reference the values data in the configmap.
	Key       string
	Name      string
	Namespace string
}

// DefaultConfig provides a default configuration to create a new ConfigMap
// Configurer by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		KubernetesClient: nil,
		Logger:           nil,
	}
}

// New creates a new configured ConfigMap Configurer.
func New(config Config) (*ConfigMapConfigurer, error) {
	if config.KubernetesClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "kubernetes client must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
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

	config.Logger.Log("debug", "checking connection to Kubernetes")
	if _, err := config.KubernetesClient.CoreV1().Namespaces().Get("default", v1.GetOptions{}); err != nil {
		return nil, microerror.MaskAny(err)
	}

	// Create a temporary file to use for holding the values file for Helm to read.
	tempFile, err := ioutil.TempFile("", "draughtsman")
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	configurer := &ConfigMapConfigurer{
		// Dependencies.
		client: config.KubernetesClient,
		logger: config.Logger,

		// Settings.
		key:       config.Key,
		name:      config.Name,
		namespace: config.Namespace,
		tempFile:  tempFile,
	}

	return configurer, nil
}

// ConfigMapConfigurer is an implementation of the Configurer interface,
// that uses a Kubernetes ConfigMap to hold configuration.
type ConfigMapConfigurer struct {
	// Dependencies.
	client kubernetes.Interface
	logger micrologger.Logger

	// Settings.
	key       string
	name      string
	namespace string
	tempFile  *os.File
}

func (c *ConfigMapConfigurer) File() (string, error) {
	defer updateConfigMapMetrics(time.Now())

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
