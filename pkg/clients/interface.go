package clients

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"kubevirt.io/client-go/kubecli"
)

type Clients interface {
	GetDynamicKubeClient() dynamic.Interface
	GetKubeVirtClient() kubecli.KubevirtClient
	GetRestClient() *rest.RESTClient
}
