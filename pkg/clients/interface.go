package clients

import (
	"k8s.io/client-go/dynamic"
)

type Clients interface {
	GetDynamicKubeClient() dynamic.Interface
}
