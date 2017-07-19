package secret

import (
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	"github.com/giantswarm/valuemodifier"
	decodemodifier "github.com/giantswarm/valuemodifier/base64/decode"

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
		return nil, microerror.MaskAnyf(invalidConfigError, "kubernetes client must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	// Settings.
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
	_, err := config.KubernetesClient.CoreV1().Namespaces().Get("default", v1.GetOptions{})
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	// The content of the secret manifests we are fetching is base64 decoded. We
	// need an decoder to write the secret data to the values file Helm can use.
	// We use the decoder from the valuemodifier package we also use to decode
	// them in opsctl.
	var decodeModifier valuemodifier.ValueModifier
	{
		modifierConfig := decodemodifier.DefaultConfig()
		decodeModifier, err = decodemodifier.New(modifierConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	configurer := &SecretConfigurer{
		// Dependencies.
		kubernetesClient: config.KubernetesClient,
		logger:           config.Logger,

		// Internals.
		decodeModifier: decodeModifier,

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

	// Internals.
	decodeModifier valuemodifier.ValueModifier

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
		return "", microerror.MaskAny(err)
	}

	var valuesData string
	{
		b, ok := s.Data[c.key]
		if !ok {
			return "", microerror.MaskAnyf(keyMissingError, "key '%v' not found in secret", c.key)
		}
		m, err := c.decodeModifier.Modify(b)
		if err != nil {
			return "", microerror.MaskAny(err)
		}
		valuesData = string(m)
	}

	return valuesData, nil
}
