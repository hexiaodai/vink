package virtualmachine

import (
	"context"
	"encoding/json"

	"fmt"

	"github.com/kubevm.io/vink/apis/annotation"
	apitypes "github.com/kubevm.io/vink/apis/types"
	"github.com/kubevm.io/vink/internal/controller/pkg"
	"github.com/kubevm.io/vink/pkg/log"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/types"
	kubevirtv1 "kubevirt.io/api/core/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DiskReconciler struct {
	Client client.Client
	Cache  cache.Cache
}

func (reconciler *DiskReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	var vm kubevirtv1.VirtualMachine
	if err := reconciler.Client.Get(ctx, request.NamespacedName, &vm); err != nil {
		if apierr.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get virtual machine: %w", err)
	}

	diskMap := make(map[string]*apitypes.VirtualMachineDisk, len(vm.Spec.Template.Spec.Volumes)-1)
	for _, disk := range vm.Spec.Template.Spec.Volumes {
		if disk.DataVolume == nil {
			continue
		}
		dv := cdiv1beta1.DataVolume{}
		if err := reconciler.Client.Get(ctx, types.NamespacedName{Namespace: vm.Namespace, Name: disk.DataVolume.Name}, &dv); err != nil {
			if apierr.IsNotFound(err) {
				continue
			}
			return ctrl.Result{}, fmt.Errorf("failed to get DataVolume: %w", err)
		}
		disk := &apitypes.VirtualMachineDisk{
			Name:             disk.Name,
			Capacity:         dv.Spec.PVC.Resources.Requests.Storage().String(),
			StorageClassName: lo.FromPtr(dv.Spec.PVC.StorageClassName),
			Mounted:          false,
			Status:           "Not Ready",
		}
		if len(dv.Spec.PVC.AccessModes) > 0 {
			disk.AccessMode = string(dv.Spec.PVC.AccessModes[0])
		}
		readyConditions := lo.Filter(dv.Status.Conditions, func(condition cdiv1beta1.DataVolumeCondition, idx int) bool {
			return condition.Type == cdiv1beta1.DataVolumeReady
		})
		if len(readyConditions) > 0 && readyConditions[0].Status == corev1.ConditionTrue {
			disk.Status = "Ready"
		}
		diskMap[disk.Name] = disk
	}

	for _, disk := range vm.Spec.Template.Spec.Domain.Devices.Disks {
		tmp, ok := diskMap[disk.Name]
		if !ok {
			continue
		}
		if lo.FromPtr(disk.BootOrder) == 1 {
			tmp.Rootfs = true
		}
		tmp.Mounted = true
		diskMap[disk.Name] = tmp
	}

	if err := pkg.PatchAnnotations(ctx, reconciler.Client, &vm, annotation.VinkDisks.Name, lo.Values(diskMap)); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (reconciler *DiskReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("virtualmachine_disk").
		For(&kubevirtv1.VirtualMachine{}).
		WatchesMetadata(
			&cdiv1beta1.DataVolume{},
			handler.TypedEnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
				metadata, ok := obj.(*metav1.PartialObjectMetadata)
				if !ok || metadata.Annotations == nil {
					return nil
				}
				ownerString, ok := metadata.Annotations[annotation.VinkDatavolumeOwner.Name]
				if !ok || len(ownerString) == 0 {
					return nil
				}
				owners := make([]string, 0)
				if err := json.Unmarshal([]byte(ownerString), &owners); err != nil {
					log.Errorf("Failed to unmarshal DataVolume owners: %v", err)
					return nil
				}
				requests := make([]reconcile.Request, 0, len(owners))
				for _, owner := range owners {
					requests = append(requests, reconcile.Request{NamespacedName: client.ObjectKey{Namespace: metadata.Namespace, Name: owner}})
				}
				return requests
			}),
		).
		Complete(reconciler)
}
