package cache

import (
	"context"

	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/informer"
)

type KubeCache struct {
	clients         clients.Clients
	InformerFactory informer.KubeInformerFactory
}

func NewKubeCache(clients clients.Clients) *KubeCache {
	kubeInformerFactory := informer.NewKubeInformerFactory(
		clients.GetVinkRestClient(),
		clients.GetKubeVirtClient().RestClient(),
		clients.GetKubeovnRestClient(),
		clients.GetKubeVirtClient(),
	)

	return &KubeCache{
		InformerFactory: kubeInformerFactory,
	}
}

func (kc *KubeCache) Start(ctx context.Context) error {
	_ = kc.InformerFactory.VirtualMachine()
	_ = kc.InformerFactory.VirtualMachineInstances()
	_ = kc.InformerFactory.DataVolume()
	_ = kc.InformerFactory.VirtualMachineSummary()
	_ = kc.InformerFactory.Subnet()

	kc.InformerFactory.Start(ctx.Done())
	kc.InformerFactory.WaitForCacheSync(ctx.Done())

	return nil
}
