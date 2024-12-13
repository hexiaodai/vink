package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kubevm.io/vink/apis/annotation"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/samber/lo"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	kubevirtv1 "kubevirt.io/api/core/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	virtualmachineFinalizer = "vink.kubevm.io/virtualmachine"
)

type DataVolumeOwnerReconciler struct {
	Client client.Client
	Cache  cache.Cache
}

func (reconciler *DataVolumeOwnerReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	var vm kubevirtv1.VirtualMachine
	if err := reconciler.Client.Get(ctx, request.NamespacedName, &vm); err != nil {
		if apierr.IsNotFound(err) {
			return ctrl.Result{}, nil
		} else {
			return ctrl.Result{}, fmt.Errorf("failed to get VirtualMachine: %w", err)
		}
	}

	if vm.DeletionTimestamp == nil && !lo.Contains(vm.Finalizers, virtualmachineFinalizer) {
		vm.Finalizers = append(vm.Finalizers, virtualmachineFinalizer)
		if err := reconciler.Client.Update(ctx, &vm); err != nil {
			if apierr.IsConflict(err) {
				return ctrl.Result{Requeue: true}, nil
			}
			return ctrl.Result{}, err
		}
	}

	for _, vol := range vm.Spec.Template.Spec.Volumes {
		if vol.DataVolume == nil {
			continue
		}

		dv := cdiv1beta1.DataVolume{}
		if err := reconciler.Client.Get(ctx, types.NamespacedName{Namespace: vm.Namespace, Name: vol.DataVolume.Name}, &dv); err != nil {
			if apierr.IsNotFound(err) {
				continue
			}
			return ctrl.Result{}, err
		}

		if dv.Annotations == nil {
			dv.Annotations = map[string]string{}
		}
		owners := make([]string, 0)
		_ = json.Unmarshal([]byte(dv.Annotations[annotation.VinkDatavolumeOwner.Name]), &owners)

		if vm.DeletionTimestamp != nil {
			owners = lo.Filter(owners, func(item string, index int) bool {
				return item != vm.Name
			})
		} else if lo.Contains(owners, vm.Name) {
			continue
		} else {
			owners = append(owners, vm.Name)
		}
		deduped := lo.SliceToMap(owners, func(item string) (string, struct{}) {
			return item, struct{}{}
		})
		owners = make([]string, 0, len(deduped))
		for key := range deduped {
			owners = append(owners, key)
		}

		var (
			bindingAnnoValue []byte
			err              error
		)
		if len(owners) > 0 {
			bindingAnnoValue, err = json.Marshal(owners)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		dv.Annotations[annotation.VinkDatavolumeOwner.Name] = string(bindingAnnoValue)
		if err := reconciler.Client.Update(ctx, &dv); err != nil {
			if apierr.IsConflict(err) {
				return ctrl.Result{Requeue: true}, nil
			}
			return ctrl.Result{}, err
		}
	}

	if vm.DeletionTimestamp != nil && lo.Contains(vm.Finalizers, virtualmachineFinalizer) {
		vm.Finalizers = lo.Filter(vm.Finalizers, func(item string, index int) bool {
			return item != virtualmachineFinalizer
		})

		if err := reconciler.Client.Update(ctx, &vm); err != nil {
			if apierr.IsConflict(err) {
				return ctrl.Result{Requeue: true}, nil
			}
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (reconciler *DataVolumeOwnerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		return syncAllVirtualMachineBindings()
	})
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named("datavolume_owner").
		For(&kubevirtv1.VirtualMachine{}).
		Complete(reconciler)
}

func syncAllVirtualMachineBindings() error {
	ctx := context.TODO()

	dvList, err := clients.Clients.CdiClient().CdiV1beta1().DataVolumes(metav1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	vmList, err := clients.Clients.VirtualMachine(metav1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	mp := make(map[types.NamespacedName]*kubevirtv1.VirtualMachine, len(dvList.Items))
	for _, vm := range vmList.Items {
		for _, vol := range vm.Spec.Template.Spec.Volumes {
			if vol.DataVolume == nil {
				continue
			}
			mp[types.NamespacedName{Namespace: vm.Namespace, Name: vol.DataVolume.Name}] = &vm
		}
	}

	bindingMap := make(map[types.NamespacedName][]string, len(dvList.Items))
	for _, dv := range dvList.Items {
		ns := types.NamespacedName{Namespace: dv.Namespace, Name: dv.Name}
		vm, ok := mp[ns]
		if !ok {
			continue
		}
		binding := bindingMap[ns]
		binding = append(binding, vm.Name)
	}

	for _, dv := range dvList.Items {
		var (
			value []byte
			err   error
		)

		ns := types.NamespacedName{Namespace: dv.Namespace, Name: dv.Name}
		binding := bindingMap[ns]
		if len(binding) == 0 && (dv.Annotations == nil || len(dv.Annotations[annotation.VinkDatavolumeOwner.Name]) == 0) {
			continue
		}

		if len(binding) > 0 {
			value, err = json.Marshal(binding)
			if err != nil {
				return err
			}
		}

		if dv.Annotations == nil {
			dv.Annotations = make(map[string]string, 1)
		}
		dv.Annotations[annotation.VinkDatavolumeOwner.Name] = string(value)
		if _, err := clients.Clients.CdiClient().CdiV1beta1().DataVolumes(dv.Namespace).Update(ctx, &dv, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}

	return nil
}
