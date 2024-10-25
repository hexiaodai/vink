package resource_event_listener

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kubevm.io/vink/apis/types"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/informer"
	"github.com/kubevm.io/vink/pkg/log"
	"k8s.io/apimachinery/pkg/runtime/schema"
	pkg_types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	client_cache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type ResourceEventListener interface {
	StartListening(ctx context.Context)
	AddSubscription(gvr schema.GroupVersionResource, crds []*types.ObjectMeta) <-chan *ResourceEvent
	RemoveSubscription(gvr schema.GroupVersionResource, eventChan <-chan *ResourceEvent)
}

func NewResourceEventListener(informerFactory informer.KubeInformerFactory) ResourceEventListener {
	r := resourceEventListener{
		informerFactory: informerFactory,
		eventListeners:  map[schema.GroupVersionResource][]chan *ResourceEvent{},
		filterFuncs:     map[<-chan *ResourceEvent][]filterFunc{},
	}

	queues := make(map[schema.GroupVersionResource]workqueue.RateLimitingInterface, len(r.informerFactory.Informers()))
	for gvr := range r.informerFactory.Informers() {
		queues[gvr] = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), gvr.String())
	}
	r.queues = queues

	return &r
}

type filterFunc func(event *ResourceEvent) (bool, error)

type resourceEventListener struct {
	informerFactory informer.KubeInformerFactory

	eventListeners map[schema.GroupVersionResource][]chan *ResourceEvent

	filterFuncs map[<-chan *ResourceEvent][]filterFunc

	queues map[schema.GroupVersionResource]workqueue.RateLimitingInterface

	mux sync.RWMutex
}

type ResourceEvent struct {
	Type    watch.EventType
	GVR     schema.GroupVersionResource
	Payload any
}

func (r *resourceEventListener) StartListening(ctx context.Context) {
	defer runtime.HandleCrash()

	defer func() {
		for _, queue := range r.queues {
			queue.ShutDown()
		}
	}()

	if err := r.registerEventHandler(); err != nil {
		panic(err)
	}

	go wait.Until(r.runWorker, 10*time.Minute, ctx.Done())
	<-ctx.Done()
}

func (r *resourceEventListener) AddSubscription(gvr schema.GroupVersionResource, metadatas []*types.ObjectMeta) <-chan *ResourceEvent {
	eventChan := make(chan *ResourceEvent, 10)

	r.mux.Lock()
	r.eventListeners[gvr] = append(r.eventListeners[gvr], eventChan)
	r.filterFuncs[eventChan] = []filterFunc{namespaceNameFilterFunc(metadatas)}
	r.mux.Unlock()

	go r.pushExistingResources(gvr, eventChan, metadatas)

	return eventChan
}

func (r *resourceEventListener) RemoveSubscription(gvr schema.GroupVersionResource, eventChan <-chan *ResourceEvent) {
	r.mux.Lock()
	defer r.mux.Unlock()

	subs, ok := r.eventListeners[gvr]
	if !ok {
		return
	}

	for i, sub := range r.eventListeners[gvr] {
		if sub != eventChan {
			continue
		}
		r.eventListeners[gvr] = append(subs[:i], subs[i+1:]...)
		delete(r.filterFuncs, eventChan)
		close(sub)
		return
	}
}

func (r *resourceEventListener) notifySubscribers(event *ResourceEvent) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	subs, ok := r.eventListeners[event.GVR]
	if !ok {
		return
	}

	for _, sub := range subs {
		fns, ok := r.filterFuncs[sub]
		if !ok {
			continue
		}
		if ok, err := r.filters(fns, event); !ok || err != nil {
			continue
		}
		sub <- event
	}
}

func (r *resourceEventListener) filters(filterFuncs []filterFunc, event *ResourceEvent) (bool, error) {
	for _, filterFunc := range filterFuncs {
		if ok, err := filterFunc(event); !ok || err != nil {
			return false, err
		}
	}
	return true, nil
}

func namespaceNameFilterFunc(items []*types.ObjectMeta) filterFunc {
	idx := make(map[string]struct{}, len(items))
	for _, metadata := range items {
		nn := pkg_types.NamespacedName{Namespace: metadata.Namespace, Name: metadata.Name}
		idx[nn.String()] = struct{}{}
	}

	return func(event *ResourceEvent) (bool, error) {
		if event.Type == watch.Deleted {
			_, ok := idx[event.Payload.(*pkg_types.NamespacedName).String()]
			return ok, nil
		}

		// unstruct, err := clients.InterfaceToUnstructured(event.Payload)
		objectMeta, err := clients.InterfaceToObjectMeta(event.Payload)
		// crd, err := utils.ConvertToCustomResourceDefinition(event.Payload)
		if err != nil {
			return false, err
		}
		nn := pkg_types.NamespacedName{Namespace: objectMeta.GetNamespace(), Name: objectMeta.GetName()}
		_, ok := idx[nn.String()]
		return ok, nil
	}
}

func resourceVersionFilterFunc(items []*types.ObjectMeta) filterFunc {
	idx := make(map[string]*types.ObjectMeta, len(items))
	for _, metadata := range items {
		ns := pkg_types.NamespacedName{Namespace: metadata.Namespace, Name: metadata.Name}
		idx[ns.String()] = metadata
	}

	return func(event *ResourceEvent) (bool, error) {
		if event.Type == watch.Deleted {
			return true, nil
		}

		objectMeta, err := clients.InterfaceToObjectMeta(event.Payload)
		// unstruct, err := clients.InterfaceToUnstructured(event.Payload)
		// crd, err := utils.ConvertToCustomResourceDefinition(event.Payload)
		if err != nil {
			return false, err
		}
		ns := pkg_types.NamespacedName{Namespace: objectMeta.GetNamespace(), Name: objectMeta.GetName()}
		original, ok := idx[ns.String()]
		if !ok {
			return false, nil
		}

		if string(objectMeta.GetUID()) != original.Uid {
			return true, err
		}

		originalVersion, err := strconv.Atoi(original.ResourceVersion)
		if err != nil {
			return false, err
		}
		version, err := strconv.Atoi(objectMeta.GetResourceVersion())
		if err != nil {
			return false, err
		}

		return version > originalVersion, nil
	}
}

func (r *resourceEventListener) pushExistingResources(gvr schema.GroupVersionResource, eventChan chan *ResourceEvent, metadatas []*types.ObjectMeta) {
	informer, ok := r.informerFactory.InformerForGVR(gvr)
	if !ok {
		return
	}

	indexer := informer.GetIndexer()
	for _, metadata := range metadatas {
		nn := pkg_types.NamespacedName{Namespace: metadata.Namespace, Name: metadata.Name}
		obj, ok, err := indexer.GetByKey(nn.String())
		if err != nil || !ok {
			continue
		}
		fns, ok := r.filterFuncs[eventChan]
		if !ok {
			// FIXME:
			continue
		}
		event := ResourceEvent{GVR: gvr, Payload: obj, Type: watch.Modified}
		if ok, err := r.filters(fns, &event); !ok || err != nil {
			continue
		}
		eventChan <- &event
	}
}

// handleResourceEvent handles a resource event from an informer and notifies subscribers.
func (r *resourceEventListener) handleResourceEvent(gvr schema.GroupVersionResource, key interface{}) error {
	informer, ok := r.informerFactory.InformerForGVR(gvr)
	if !ok {
		return fmt.Errorf("informer not found for %s", gvr)
	}

	event := ResourceEvent{GVR: gvr}

	obj, exists, err := informer.GetIndexer().GetByKey(key.(string))
	if err != nil {
		return err
	}
	if exists {
		event.Type = watch.Modified
		event.Payload = obj
	} else {
		nn := pkg_types.NamespacedName{}
		parts := strings.Split(key.(string), "/")
		switch len(parts) {
		case 1:
			nn.Name = parts[0]
		case 2:
			nn.Namespace = parts[0]
			nn.Name = parts[1]
		}
		event.Type = watch.Deleted
		event.Payload = &nn
	}

	r.notifySubscribers(&event)

	return nil
}

func (r *resourceEventListener) registerEventHandler() error {
	for gvr, informer := range r.informerFactory.Informers() {
		eventHandler := client_cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, err := client_cache.MetaNamespaceKeyFunc(obj)
				if err != nil {
					log.Errorf("failed to get key for %v when adding object: %v", gvr, err)
					return
				}
				if q, ok := r.queues[gvr]; ok {
					q.Add(key)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				key, err := client_cache.MetaNamespaceKeyFunc(newObj)
				if err != nil {
					log.Errorf("failed to get key for %v when updating object: %v", gvr, err)
					return
				}
				if q, ok := r.queues[gvr]; ok {
					q.Add(key)
				}
			},
			DeleteFunc: func(obj interface{}) {
				key, err := client_cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
				if err != nil {
					log.Errorf("failed to get key for %v when deleting object: %v", gvr, err)
					return
				}
				if q, ok := r.queues[gvr]; ok {
					q.Add(key)
				}
			},
		}

		log.Info(fmt.Sprintf("adding event handler for %v", gvr.String()))
		if _, err := informer.AddEventHandler(eventHandler); err != nil {
			return fmt.Errorf("failed to add event handler for GVR: %v: %v", gvr, err)
		}
	}

	return nil
}

func (r *resourceEventListener) runWorker() {
	wg := sync.WaitGroup{}
	wg.Add(len(r.queues))

	for gvr := range r.queues {
		go func(gvr schema.GroupVersionResource) {
			defer wg.Done()
			for r.processNextItem(gvr) {
			}
		}(gvr)
	}
	wg.Wait()
}

func (r *resourceEventListener) processNextItem(gvr schema.GroupVersionResource) bool {
	queue, ok := r.queues[gvr]
	if !ok {
		return false
	}

	key, quit := queue.Get()
	if quit {
		return false
	}
	defer queue.Done(key)

	err := r.handleResourceEvent(gvr, key)
	r.handleErr(err, gvr, key)

	return true
}

func (r *resourceEventListener) handleErr(err error, gvr schema.GroupVersionResource, key interface{}) {
	queue, ok := r.queues[gvr]
	if !ok {
		return
	}

	if err == nil {
		queue.Forget(key)
		return
	}

	if queue.NumRequeues(key) < 5 {
		log.Infof("error syncing '%v' %v: %v", gvr.String(), key, err)
		queue.AddRateLimited(key)
		return
	}

	queue.Forget(key)
	runtime.HandleError(err)
	log.Infof("dropping %v %q out of the queue: %v", gvr.String(), key, err)
}
