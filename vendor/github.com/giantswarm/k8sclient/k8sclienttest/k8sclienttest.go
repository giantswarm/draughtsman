package k8sclienttest

import (
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/k8sclient/k8scrdclient"
)

type ClientsConfig struct {
	CrdClient  k8scrdclient.Interface
	CtrlClient client.Client
	DynClient  dynamic.Interface
	ExtClient  apiextensionsclient.Interface
	G8sClient  versioned.Interface
	K8sClient  kubernetes.Interface
	RestClient rest.Interface
	RestConfig *rest.Config
}

type Clients struct {
	crdClient  k8scrdclient.Interface
	ctrlClient client.Client
	dynClient  dynamic.Interface
	extClient  apiextensionsclient.Interface
	g8sClient  versioned.Interface
	k8sClient  kubernetes.Interface
	restClient rest.Interface
	restConfig *rest.Config
}

func NewClients(config ClientsConfig) (*Clients, error) {
	c := &Clients{
		crdClient:  config.CrdClient,
		ctrlClient: config.CtrlClient,
		dynClient:  config.DynClient,
		extClient:  config.ExtClient,
		g8sClient:  config.G8sClient,
		k8sClient:  config.K8sClient,
		restClient: config.RestClient,
		restConfig: config.RestConfig,
	}

	return c, nil
}

func (c *Clients) CRDClient() k8scrdclient.Interface {
	return c.crdClient
}

func (c *Clients) CtrlClient() client.Client {
	return c.ctrlClient
}

func (c *Clients) DynClient() dynamic.Interface {
	return c.dynClient
}

func (c *Clients) ExtClient() apiextensionsclient.Interface {
	return c.extClient
}

func (c *Clients) G8sClient() versioned.Interface {
	return c.g8sClient
}

func (c *Clients) K8sClient() kubernetes.Interface {
	return c.k8sClient
}

func (c *Clients) RESTClient() rest.Interface {
	return c.restClient
}

func (c *Clients) RESTConfig() *rest.Config {
	return rest.CopyConfig(c.restConfig)
}

func (c *Clients) Scheme() *runtime.Scheme {
	return scheme.Scheme
}
