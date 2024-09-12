package clients

import (
	"k8s.io/client-go/dynamic"
	"kubevirt.io/client-go/kubecli"
)

type Clients interface {
	GetDynamicKubeClient() dynamic.Interface
	GetKubeVirtClient() kubecli.KubevirtClient
}
