package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/kubevm.io/vink/pkg/log"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/util/workqueue"
)

type Reconciler interface {
	Reconcile(ctx context.Context, key string) error
}

func NewController(informer cache.SharedIndexInformer, reconciler Reconciler) *Controller {
	return &Controller{
		informer:   informer,
		reconciler: reconciler,
		queue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), ""),
	}
}

type Controller struct {
	reconciler Reconciler
	queue      workqueue.RateLimitingInterface
	informer   cache.SharedIndexInformer
}

func (c *Controller) Start(ctx context.Context) {
	defer runtime.HandleCrash()

	defer func() {
		c.queue.ShutDown()
	}()

	if err := c.eventHandler(); err != nil {
		panic(err)
	}

	go wait.Until(c.runWorker, time.Second, ctx.Done())
	<-ctx.Done()
}

func (c *Controller) eventHandler() error {
	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    c.enqueueObject,
		UpdateFunc: func(_, newObj interface{}) { c.enqueueObject(newObj) },
		DeleteFunc: c.enqueueDeletion,
	}

	if _, err := c.informer.AddEventHandler(eventHandler); err != nil {
		return fmt.Errorf("failed to add event handler: %v", err)
	}

	return nil
}

func (c *Controller) enqueueObject(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Errorf("failed to get key for object when adding: %w", err)
		return
	}
	c.queue.Add(key)
}

func (c *Controller) enqueueDeletion(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Errorf("failed to get key for object when deleting: %w", err)
		return
	}
	c.queue.Add(key)
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	ctx := context.TODO()
	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		return c.reconciler.Reconcile(ctx, key.(string))
	})

	c.handleErr(err, key)

	return true
}

func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if k8serr.IsConflict(err) {
		c.queue.AddRateLimited(key)
		return
	}

	if c.queue.NumRequeues(key) < 5 {
		log.Errorf("error syncing '%v': %v", key, err)
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	runtime.HandleError(err)
	log.Infof("dropping '%v' out of the queue: %v", key, err)
}
