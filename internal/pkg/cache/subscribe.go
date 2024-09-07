package cache

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kubevm.io/vink/apis/types"
	"github.com/kubevm.io/vink/pkg/log"
	"github.com/kubevm.io/vink/pkg/utils"
	"k8s.io/apimachinery/pkg/runtime/schema"
	pkg_types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	client_cache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Subscribe interface {
	Run(ctx context.Context)
	Subscribe(gvr schema.GroupVersionResource, crds []*types.ObjectMeta) <-chan *ResourceEvent
	Unsubscribe(gvr schema.GroupVersionResource, eventChan <-chan *ResourceEvent)
}

func NewSubscribe(cache Cache) Subscribe {
	s := subscribe{
		cache:       cache,
		subscribers: map[schema.GroupVersionResource][]chan *ResourceEvent{},
		filterFuncs: sync.Map{},
	}

	informers := s.cache.ListInformers()
	queue := make(map[schema.GroupVersionResource]workqueue.RateLimitingInterface, len(informers))

	for gvr := range informers {
		queue[gvr] = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), gvr.String())
	}
	s.queue = queue

	return &s
}

type filterFunc func(event *ResourceEvent) (bool, error)

type subscribe struct {
	cache Cache

	subscribers map[schema.GroupVersionResource][]chan *ResourceEvent

	// key: <-chan *ResourceEvent
	// value: []filterFunc
	filterFuncs sync.Map

	queue map[schema.GroupVersionResource]workqueue.RateLimitingInterface

	mux sync.RWMutex
}

type ResourceEvent struct {
	Type    watch.EventType
	GVR     schema.GroupVersionResource
	Payload interface{}
}

func (s *subscribe) Run(ctx context.Context) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer func() {
		for _, queue := range s.queue {
			queue.ShutDown()
		}
	}()

	if err := s.registerEventHandler(); err != nil {
		panic(err)
	}

	go wait.Until(s.runWorker, 10*time.Minute, ctx.Done())
	<-ctx.Done()
}

// subscribe allows a client to subscribe to a GVR resource's events.
func (s *subscribe) Subscribe(gvr schema.GroupVersionResource, metadatas []*types.ObjectMeta) <-chan *ResourceEvent {
	// Buffer size can be adjusted
	eventChan := make(chan *ResourceEvent, 10)

	// Add the new subscriber to the list
	s.mux.Lock()
	s.subscribers[gvr] = append(s.subscribers[gvr], eventChan)
	s.filterFuncs.Store(eventChan, []filterFunc{namespaceNameFilterFunc(metadatas), resourceVersionFilterFunc(metadatas)})
	s.mux.Unlock()

	// Push existing resources to the new subscriber
	go s.pushExistingResources(gvr, eventChan, metadatas)

	return eventChan
}

func (s *subscribe) Unsubscribe(gvr schema.GroupVersionResource, eventChan <-chan *ResourceEvent) {
	s.mux.Lock()
	defer s.mux.Unlock()

	subs, ok := s.subscribers[gvr]
	if !ok {
		return
	}

	for i, sub := range s.subscribers[gvr] {
		if sub != eventChan {
			continue
		}
		s.subscribers[gvr] = append(subs[:i], subs[i+1:]...)
		s.filterFuncs.Delete(sub)
		close(sub)
		return
	}
}

// notifySubscribers notifies all subscribers of a resource event.
func (s *subscribe) notifySubscribers(event *ResourceEvent) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	subs, ok := s.subscribers[event.GVR]
	if !ok {
		return
	}

	// Send the event to all subscribers
	for _, sub := range subs {
		value, ok := s.filterFuncs.Load(sub)
		if !ok {
			continue
		}
		if ok, err := s.filters(value.([]filterFunc), event); !ok || err != nil {
			continue
		}
		sub <- event
	}
}

func (s *subscribe) filters(filterFuncs []filterFunc, event *ResourceEvent) (bool, error) {
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

		crd, err := utils.ConvertToCustomResourceDefinition(event.Payload)
		if err != nil {
			return false, err
		}
		nn := pkg_types.NamespacedName{Namespace: crd.Metadata.Namespace, Name: crd.Metadata.Name}
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

		crd, err := utils.ConvertToCustomResourceDefinition(event.Payload)
		if err != nil {
			return false, err
		}
		ns := pkg_types.NamespacedName{Namespace: crd.Metadata.Namespace, Name: crd.Metadata.Name}
		original, ok := idx[ns.String()]
		if !ok {
			return false, nil
		}

		if crd.Metadata.Uid != original.Uid {
			return true, err
		}

		originalVersion, err := strconv.Atoi(original.ResourceVersion)
		if err != nil {
			return false, err
		}
		version, err := strconv.Atoi(crd.Metadata.ResourceVersion)
		if err != nil {
			return false, err
		}

		return version > originalVersion, nil
	}
}

func (s *subscribe) pushExistingResources(gvr schema.GroupVersionResource, eventChan chan<- *ResourceEvent, metadatas []*types.ObjectMeta) {
	informer, err := s.cache.GetInformer(gvr)
	if err != nil {
		return
	}
	indexer := informer.GetIndexer()
	for _, metadata := range metadatas {
		nn := pkg_types.NamespacedName{Namespace: metadata.Namespace, Name: metadata.Name}
		obj, ok, err := indexer.GetByKey(nn.String())
		if err != nil || !ok {
			continue
		}
		value, ok := s.filterFuncs.Load(eventChan)
		if !ok {
			// FIXME:
			continue
		}
		event := ResourceEvent{GVR: gvr, Payload: obj, Type: watch.Modified}
		if ok, err := s.filters(value.([]filterFunc), &event); !ok || err != nil {
			continue
		}
		eventChan <- &event
	}
}

// handleResourceEvent handles a resource event from an informer and notifies subscribers.
func (s *subscribe) handleResourceEvent(gvr schema.GroupVersionResource, key interface{}) error {
	informer, err := s.cache.GetInformer(gvr)
	if err != nil {
		return err
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

	s.notifySubscribers(&event)

	return nil
}

func (s *subscribe) registerEventHandler() error {
	for gvr, informer := range s.cache.ListInformers() {
		eventHandler := client_cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, err := client_cache.MetaNamespaceKeyFunc(obj)
				if err != nil {
					log.Errorf("failed to get key for %v when adding object: %v", gvr, err)
					return
				}
				if q, ok := s.queue[gvr]; ok {
					q.Add(key)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				key, err := client_cache.MetaNamespaceKeyFunc(newObj)
				if err != nil {
					log.Errorf("failed to get key for %v when updating object: %v", gvr, err)
					return
				}
				if q, ok := s.queue[gvr]; ok {
					q.Add(key)
				}
			},
			DeleteFunc: func(obj interface{}) {
				key, err := client_cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
				if err != nil {
					log.Errorf("failed to get key for %v when deleting object: %v", gvr, err)
					return
				}
				if q, ok := s.queue[gvr]; ok {
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

func (s *subscribe) runWorker() {
	wg := sync.WaitGroup{}
	wg.Add(len(s.queue))

	for gvr := range s.queue {
		go func(gvr schema.GroupVersionResource) {
			defer wg.Done()
			for s.processNextItem(gvr) {
			}
		}(gvr)
	}
	wg.Wait()
}

func (s *subscribe) processNextItem(gvr schema.GroupVersionResource) bool {
	queue, ok := s.queue[gvr]
	if !ok {
		return false
	}

	key, quit := queue.Get()
	if quit {
		return false
	}
	defer queue.Done(key)

	err := s.handleResourceEvent(gvr, key)
	s.handleErr(err, gvr, key)

	return true
}

func (s *subscribe) handleErr(err error, gvr schema.GroupVersionResource, key interface{}) {
	queue, ok := s.queue[gvr]
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
