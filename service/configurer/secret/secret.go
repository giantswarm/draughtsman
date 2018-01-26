package secret

import (
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/draughtsman/service/configurer/spec"
)

// ConfigurerType is the kind of a Configurer that is backed by a Kubernetes
// Secret.
var ConfigurerType spec.ConfigurerType = "SecretConfigurer"

// Config represents the configuration used to create a Secret Configurer.
type Config struct {
	// Dependencies.
	KubernetesClient kubernetes.Interface
	Logger           micrologger.Logger

	// Settings.

	// Key is the key to reference the values data in the secret.
	Key       string
	Name      string
	Namespace string
}

// DefaultConfig provides a default configuration to create a new Secret
// Configurer by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		KubernetesClient: nil,
		Logger:           nil,

		// Settings.
		Key:       "",
		Name:      "",
		Namespace: "",
	}
}

// New creates a new configured Secret Configurer.
func New(config Config) (*SecretConfigurer, error) {
	// Dependencies.
	if config.KubernetesClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "kubernetes client must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "logger must not be empty")
	}

	// Settings.
	if config.Key == "" {
		return nil, microerror.Maskf(invalidConfigError, "key must not be empty")
	}
	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "name must not be empty")
	}
	if config.Namespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "namespace must not be empty")
	}

	config.Logger.Log("debug", "checking connection to Kubernetes")
	_, err := config.KubernetesClient.CoreV1().Secrets("draughtsman").Get("draughtsman", v1.GetOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	configurer := &SecretConfigurer{
		// Dependencies.
		kubernetesClient: config.KubernetesClient,
		logger:           config.Logger,

		// Settings.
		key:       config.Key,
		name:      config.Name,
		namespace: config.Namespace,
	}

	return configurer, nil
}

// SecretConfigurer is an implementation of the Configurer interface, that uses
// a Kubernetes Secret to hold configuration.
type SecretConfigurer struct {
	// Dependencies.
	kubernetesClient kubernetes.Interface
	logger           micrologger.Logger

	// Settings.
	key       string
	name      string
	namespace string
}

func (c *SecretConfigurer) Type() spec.ConfigurerType {
	return ConfigurerType
}

func (c *SecretConfigurer) Values() (string, error) {
	defer updateSecretMetrics(time.Now())

	c.logger.Log("debug", "fetching configuration from secret", "name", c.name, "namespace", c.namespace)

	s, err := c.kubernetesClient.CoreV1().Secrets(c.namespace).Get(c.name, v1.GetOptions{})
	if err != nil {
		return "", microerror.Mask(err)
	}

	var valuesData string
	{
		b, ok := s.Data[c.key]
		if !ok {
			return "", microerror.Maskf(keyMissingError, "key '%v' not found in secret", c.key)
		}
		valuesData = string(b)
	}

	return valuesData, nil
}
