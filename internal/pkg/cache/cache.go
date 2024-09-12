package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/clients/gvr"
	"github.com/kubevm.io/vink/pkg/informer"
	"k8s.io/apimachinery/pkg/runtime/schema"
	client_cache "k8s.io/client-go/tools/cache"
	virtv1 "kubevirt.io/api/core/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

type Cache interface {
	GetInformer(gvr schema.GroupVersionResource) (client_cache.SharedIndexInformer, error)
	ListInformers() map[schema.GroupVersionResource]client_cache.SharedIndexInformer
}

type cache struct {
	clients clients.Clients

	informerFactory informer.KubeInformerFactory

	informers map[schema.GroupVersionResource]client_cache.SharedIndexInformer
	infmux    sync.RWMutex
}

func NewCache(ctx context.Context) Cache {
	c := cache{}

	c.clients = clients.GetClients()
	c.informerFactory = informer.NewKubeInformerFactory(c.clients.GetKubeVirtClient().RestClient(), c.clients.GetKubeVirtClient())

	c.informers = make(map[schema.GroupVersionResource]client_cache.SharedIndexInformer)
	c.informers[gvr.From(virtv1.VirtualMachine{})] = c.informerFactory.VirtualMachine()
	c.informers[gvr.From(virtv1.VirtualMachineInstance{})] = c.informerFactory.VirtualMachineInstances()
	c.informers[gvr.From(cdiv1.DataVolume{})] = c.informerFactory.DataVolume()

	c.informerFactory.Start(ctx.Done())
	c.informerFactory.WaitForCacheSync(ctx.Done())

	return &c
}

func (c *cache) GetInformer(gvr schema.GroupVersionResource) (client_cache.SharedIndexInformer, error) {
	c.infmux.RLock()
	defer c.infmux.RUnlock()

	informer, ok := c.informers[gvr]
	if !ok {
		return nil, fmt.Errorf("informer not found for %s", gvr)
	}
	return informer, nil
}

func (c *cache) ListInformers() map[schema.GroupVersionResource]client_cache.SharedIndexInformer {
	c.infmux.RLock()
	defer c.infmux.RUnlock()

	informersCopy := make(map[schema.GroupVersionResource]client_cache.SharedIndexInformer, len(c.informers))
	for key, value := range c.informers {
		informersCopy[key] = value
	}
	return informersCopy
}
