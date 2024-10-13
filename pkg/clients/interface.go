package clients

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"kubevirt.io/client-go/kubecli"
)

type Clients interface {
	GetDynamicKubeClient() dynamic.Interface
	GetKubeVirtClient() kubecli.KubevirtClient
	GetDiscoveryClient() discovery.DiscoveryInterface
	GetKubeConfig() *rest.Config
	GetVinkRestClient() *rest.RESTClient
	GetKubeovnRestClient() *rest.RESTClient
}
