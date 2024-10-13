package virtualmachine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kubevm.io/vink/apis/annotation"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/clients/gvr"
	"github.com/kubevm.io/vink/pkg/log"
	"github.com/samber/lo"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/util/workqueue"
	kubev1 "kubevirt.io/api/core/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

const (
	virtualmachineFinalizer = "vink.kubevm.io/virtualmachine"
)

type Controller struct {
	queue      workqueue.RateLimitingInterface
	vmInformer cache.SharedIndexInformer
}

func New(vmInformer cache.SharedIndexInformer) *Controller {
	return &Controller{
		vmInformer: vmInformer,
		queue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), gvr.From(kubev1.VirtualMachine{}).String()),
	}
}

func (c *Controller) Run(ctx context.Context) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer func() {
		c.queue.ShutDown()
	}()

	if err := c.eventHandler(); err != nil {
		panic(err)
	}

	go wait.Until(c.runWorker, 10*time.Minute, ctx.Done())
	<-ctx.Done()
}

func (c *Controller) eventHandler() error {
	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				log.Errorf("failed to get key for %v when adding object: %v", obj, err)
				return
			}
			c.queue.Add(key)
		},
		UpdateFunc: func(_, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err != nil {
				log.Errorf("failed to get key for %v when adding object: %v", newObj, err)
				return
			}
			c.queue.Add(key)
		},
		DeleteFunc: func(obj interface{}) {

		},
	}
	if _, err := c.vmInformer.AddEventHandler(eventHandler); err != nil {
		return fmt.Errorf("failed to AddEventHandler, error: %v", err)
	}

	return nil
}

func (c *Controller) setupDataVolume(key interface{}) error {
	ctx := context.TODO()

	obj, exists, err := c.vmInformer.GetIndexer().GetByKey(key.(string))
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	vm, ok := obj.(*kubev1.VirtualMachine)
	if !ok {
		return fmt.Errorf("failed to get '%v' vm object", key)
	}

	cli := clients.GetClients().GetDynamicKubeClient()

	vmCli := cli.Resource(gvr.From(kubev1.VirtualMachine{})).Namespace(vm.Namespace)
	if vm.DeletionTimestamp == nil && !lo.Contains(vm.Finalizers, virtualmachineFinalizer) {
		vm.Finalizers = append(vm.Finalizers, virtualmachineFinalizer)
		unObj, err := clients.Unstructured(vm)
		if err != nil {
			return err
		}
		value, err := vmCli.Update(ctx, unObj, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		newvm, err := clients.FromUnstructured[kubev1.VirtualMachine](value)
		if err != nil {
			return err
		}
		vm = newvm
	}

	dvCli := cli.Resource(gvr.From(cdiv1beta1.DataVolume{})).Namespace(vm.Namespace)
	for _, vol := range vm.Spec.Template.Spec.Volumes {
		if vol.DataVolume == nil {
			continue
		}
		unObj, err := dvCli.Get(ctx, vol.DataVolume.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		dvAnno := unObj.GetAnnotations()
		binding := make([]string, 0)
		_ = json.Unmarshal([]byte(dvAnno[annotation.VinkVirtualmachineBinding.Name]), &binding)

		if vm.DeletionTimestamp != nil {
			binding = lo.Filter(binding, func(item string, index int) bool {
				return item != vm.Name
			})
		} else if lo.Contains(binding, vm.Name) {
			continue
		} else {
			binding = append(binding, vm.Name)
		}
		deduped := lo.SliceToMap(binding, func(item string) (string, struct{}) {
			return item, struct{}{}
		})
		binding = make([]string, 0, len(deduped))
		for key := range deduped {
			binding = append(binding, key)
		}

		bindingAnnoValue, err := json.Marshal(binding)
		if err != nil {
			return err
		}

		unObj.SetAnnotations(map[string]string{
			annotation.VinkVirtualmachineBinding.Name: string(bindingAnnoValue),
		})
		if _, err := dvCli.Update(ctx, unObj, metav1.UpdateOptions{}); err != nil && !k8serr.IsNotFound(err) {
			return err
		}
	}

	if vm.DeletionTimestamp != nil && lo.Contains(vm.Finalizers, virtualmachineFinalizer) {
		vm.Finalizers = lo.Filter(vm.Finalizers, func(item string, index int) bool {
			return item != virtualmachineFinalizer
		})
		unObj, err := clients.Unstructured(vm)
		if err != nil {
			return err
		}
		if _, err := vmCli.Update(ctx, unObj, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}

	return nil
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

	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		return c.setupDataVolume(key)
	})
	c.handleErr(err, key)

	return true
}

func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < 5 {
		log.Errorf("error syncing '%v' %v: %v", gvr.From(kubev1.VirtualMachine{}).String(), key, err)
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	runtime.HandleError(err)
	log.Infof("dropping %v %q out of the queue: %v", gvr.From(kubev1.VirtualMachine{}).String(), key, err)
}
