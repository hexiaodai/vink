package informer

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	kubeovn "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/clients/gvr"
	"github.com/kubevm.io/vink/pkg/k8s/apis/vink/v1alpha1"
	"github.com/kubevm.io/vink/pkg/log"
	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	kubev1 "kubevirt.io/api/core/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

var unexpectedObjectError = errors.New("unexpected object")

type newSharedInformer func() cache.SharedIndexInformer

type KubeInformerFactory interface {
	// Starts any informers that have not been started yet
	// This function is thread safe and idempotent
	Start(stopCh <-chan struct{})

	// Waits for all informers to sync
	WaitForCacheSync(stopCh <-chan struct{})

	VirtualMachineInstances() cache.SharedIndexInformer

	// VirtualMachine handles the VMIs that are stopped or not running
	VirtualMachine() cache.SharedIndexInformer

	VirtualMachineSummary() cache.SharedIndexInformer

	DataVolume() cache.SharedIndexInformer

	Node() cache.SharedIndexInformer

	Subnet() cache.SharedIndexInformer

	Event() cache.SharedIndexInformer

	K8SInformerFactory() informers.SharedInformerFactory

	InformerForGVR(gvr schema.GroupVersionResource) (cache.SharedIndexInformer, bool)

	Informers() map[schema.GroupVersionResource]cache.SharedIndexInformer
}

type kubeInformerFactory struct {
	lock          sync.Mutex
	defaultResync time.Duration

	informers        map[schema.GroupVersionResource]cache.SharedIndexInformer
	startedInformers map[schema.GroupVersionResource]bool
	k8sInformers     informers.SharedInformerFactory
}

func NewKubeInformerFactory() KubeInformerFactory {
	return &kubeInformerFactory{
		// Resulting resync period will be between 12 and 24 hours, like the default for k8s
		defaultResync:    resyncPeriod(12 * time.Hour),
		informers:        make(map[schema.GroupVersionResource]cache.SharedIndexInformer),
		startedInformers: make(map[schema.GroupVersionResource]bool),
		k8sInformers:     informers.NewSharedInformerFactoryWithOptions(clients.Instance.KubevirtClient, 0),
	}
}

// Start can be called from multiple controllers in different go routines safely.
// Only informers that have not started are triggered by this function.
// Multiple calls to this function are idempotent.
func (f *kubeInformerFactory) Start(stopCh <-chan struct{}) {
	f.lock.Lock()
	defer f.lock.Unlock()

	for name, informer := range f.informers {
		if f.startedInformers[name] {
			// skip informers that have already started.
			log.Debugf("SKIPPING informer %T", name)
			continue
		}
		log.Infof("STARTING informer %T", name)
		go informer.Run(stopCh)
		f.startedInformers[name] = true
	}
	f.k8sInformers.Start(stopCh)
}

func (f *kubeInformerFactory) WaitForCacheSync(stopCh <-chan struct{}) {
	syncs := []cache.InformerSynced{}

	f.lock.Lock()
	for name, informer := range f.informers {
		log.Infof("Waiting for cache sync of informer %T", name)
		syncs = append(syncs, informer.HasSynced)
	}
	f.lock.Unlock()

	cache.WaitForCacheSync(stopCh, syncs...)
}

func GetVMIInformerIndexers() cache.Indexers {
	return cache.Indexers{
		cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
		"node": func(obj interface{}) (strings []string, e error) {
			return []string{obj.(*kubev1.VirtualMachineInstance).Status.NodeName}, nil
		},
		"dv": func(obj interface{}) ([]string, error) {
			vmi, ok := obj.(*kubev1.VirtualMachineInstance)
			if !ok {
				return nil, unexpectedObjectError
			}
			var dvs []string
			for _, vol := range vmi.Spec.Volumes {
				if vol.DataVolume != nil {
					dvs = append(dvs, fmt.Sprintf("%s/%s", vmi.Namespace, vol.DataVolume.Name))
				}
			}
			return dvs, nil
		},
		"pvc": func(obj interface{}) ([]string, error) {
			vmi, ok := obj.(*kubev1.VirtualMachineInstance)
			if !ok {
				return nil, unexpectedObjectError
			}
			var pvcs []string
			for _, vol := range vmi.Spec.Volumes {
				if vol.PersistentVolumeClaim != nil {
					pvcs = append(pvcs, fmt.Sprintf("%s/%s", vmi.Namespace, vol.PersistentVolumeClaim.ClaimName))
				}
			}
			return pvcs, nil
		},
	}
}

func (f *kubeInformerFactory) VirtualMachineInstances() cache.SharedIndexInformer {
	return f.getInformer(gvr.From(kubev1.VirtualMachineInstance{}), func() cache.SharedIndexInformer {
		lw := cache.NewListWatchFromClient(clients.Instance.KubevirtClient.RestClient(), "virtualmachineinstances", k8sv1.NamespaceAll, fields.Everything())
		return cache.NewSharedIndexInformer(lw, &kubev1.VirtualMachineInstance{}, f.defaultResync, GetVMIInformerIndexers())
	})
}

func (f *kubeInformerFactory) Subnet() cache.SharedIndexInformer {
	return f.getInformer(gvr.From(kubeovn.Subnet{}), func() cache.SharedIndexInformer {
		lw := cache.NewListWatchFromClient(clients.Instance.KubeOVNRestClient, "subnets", k8sv1.NamespaceAll, fields.Everything())
		return cache.NewSharedIndexInformer(lw, &kubeovn.Subnet{}, f.defaultResync, cache.Indexers{})
	})
}

func GetVirtualMachineInformerIndexers() cache.Indexers {
	return cache.Indexers{
		cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
		"dv": func(obj interface{}) ([]string, error) {
			vm, ok := obj.(*kubev1.VirtualMachine)
			if !ok {
				return nil, unexpectedObjectError
			}
			var dvs []string
			for _, vol := range vm.Spec.Template.Spec.Volumes {
				if vol.DataVolume != nil {
					dvs = append(dvs, fmt.Sprintf("%s/%s", vm.Namespace, vol.DataVolume.Name))
				}
			}
			return dvs, nil
		},
		"pvc": func(obj interface{}) ([]string, error) {
			vm, ok := obj.(*kubev1.VirtualMachine)
			if !ok {
				return nil, unexpectedObjectError
			}
			var pvcs []string
			for _, vol := range vm.Spec.Template.Spec.Volumes {
				if vol.PersistentVolumeClaim != nil {
					pvcs = append(pvcs, fmt.Sprintf("%s/%s", vm.Namespace, vol.PersistentVolumeClaim.ClaimName))
				}
			}
			return pvcs, nil
		},
	}
}

func (f *kubeInformerFactory) VirtualMachine() cache.SharedIndexInformer {
	return f.getInformer(gvr.From(kubev1.VirtualMachine{}), func() cache.SharedIndexInformer {
		lw := cache.NewListWatchFromClient(clients.Instance.KubevirtClient.RestClient(), "virtualmachines", k8sv1.NamespaceAll, fields.Everything())
		return cache.NewSharedIndexInformer(lw, &kubev1.VirtualMachine{}, f.defaultResync, cache.Indexers{})
	})
}

func (f *kubeInformerFactory) VirtualMachineSummary() cache.SharedIndexInformer {
	return f.getInformer(gvr.From(v1alpha1.VirtualMachineSummary{}), func() cache.SharedIndexInformer {
		lw := cache.NewListWatchFromClient(clients.Instance.VinkRestClient, "virtualmachinesummarys", k8sv1.NamespaceAll, fields.Everything())
		return cache.NewSharedIndexInformer(lw, &v1alpha1.VirtualMachineSummary{}, f.defaultResync, cache.Indexers{})
	})
}

func (f *kubeInformerFactory) DataVolume() cache.SharedIndexInformer {
	return f.getInformer(gvr.From(cdiv1.DataVolume{}), func() cache.SharedIndexInformer {
		lw := cache.NewListWatchFromClient(clients.Instance.CdiClient().CdiV1beta1().RESTClient(), "datavolumes", k8sv1.NamespaceAll, fields.Everything())
		return cache.NewSharedIndexInformer(lw, &cdiv1.DataVolume{}, f.defaultResync, cache.Indexers{})
	})
}

func (f *kubeInformerFactory) Node() cache.SharedIndexInformer {
	return f.getInformer(gvr.From(k8sv1.Node{}), func() cache.SharedIndexInformer {
		lw := cache.NewListWatchFromClient(clients.Instance.KubeRestClient, "nodes", k8sv1.NamespaceAll, fields.Everything())
		return cache.NewSharedIndexInformer(lw, &k8sv1.Node{}, f.defaultResync, cache.Indexers{})
	})
}

func (f *kubeInformerFactory) Event() cache.SharedIndexInformer {
	return f.getInformer(gvr.From(k8sv1.Event{}), func() cache.SharedIndexInformer {
		lw := cache.NewListWatchFromClient(clients.Instance.KubeRestClient, "events", k8sv1.NamespaceAll, fields.Everything())
		return cache.NewSharedIndexInformer(lw, &k8sv1.Event{}, f.defaultResync, cache.Indexers{})
	})
}

func (f *kubeInformerFactory) K8SInformerFactory() informers.SharedInformerFactory {
	return f.k8sInformers
}

func (f *kubeInformerFactory) InformerForGVR(gvr schema.GroupVersionResource) (cache.SharedIndexInformer, bool) {
	f.lock.Lock()
	defer f.lock.Unlock()

	informer, ok := f.informers[gvr]
	return informer, ok
}

func (f *kubeInformerFactory) Informers() map[schema.GroupVersionResource]cache.SharedIndexInformer {
	f.lock.Lock()
	defer f.lock.Unlock()

	informers := make(map[schema.GroupVersionResource]cache.SharedIndexInformer, len(f.informers))
	for gvr, informer := range f.informers {
		informers[gvr] = informer
	}
	return informers
}

// internal function used to retrieve an already created informer
// or create a new informer if one does not already exist.
// Thread safe
func (f *kubeInformerFactory) getInformer(key schema.GroupVersionResource, newFunc newSharedInformer) cache.SharedIndexInformer {
	f.lock.Lock()
	defer f.lock.Unlock()

	informer, exists := f.informers[key]
	if exists {
		return informer
	}
	informer = newFunc()
	f.informers[key] = informer

	return informer
}

// resyncPeriod computes the time interval a shared informer waits before resyncing with the api server
func resyncPeriod(minResyncPeriod time.Duration) time.Duration {
	// #nosec no need for better randomness
	factor := rand.Float64() + 1
	return time.Duration(float64(minResyncPeriod.Nanoseconds()) * factor)
}
