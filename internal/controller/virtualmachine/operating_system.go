package virtualmachine

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kubevm.io/vink/apis/annotation"
	"github.com/kubevm.io/vink/apis/label"
	apitypes "github.com/kubevm.io/vink/apis/types"
	"github.com/kubevm.io/vink/internal/controller/pkg"
	"github.com/kubevm.io/vink/pkg/log"
	"github.com/samber/lo"
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

type OperatingSystemReconciler struct {
	Client client.Client
	Cache  cache.Cache
}

func (reconciler *OperatingSystemReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	var vm kubevirtv1.VirtualMachine
	if err := reconciler.Client.Get(ctx, request.NamespacedName, &vm); err != nil {
		if apierr.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get VirtualMachine: %w", err)
	}

	// if vm.Annotations != nil && len(vm.Annotations[annotation.VinkOperatingSystem.Name]) > 0 {
	// 	return ctrl.Result{}, nil
	// }

	if len(vm.Spec.Template.Spec.Domain.Devices.Disks) == 0 {
		if err := pkg.PatchAnnotations(ctx, reconciler.Client, &vm, annotation.VinkOperatingSystem.Name, struct{}{}); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to patch VirtualMachine operating system: %w", err)
		}
		return ctrl.Result{}, nil
	}

	rootDiskName := vm.Spec.Template.Spec.Domain.Devices.Disks[0].Name

	bootDisks := lo.Filter(vm.Spec.Template.Spec.Domain.Devices.Disks, func(disk kubevirtv1.Disk, idx int) bool {
		return lo.FromPtr(disk.BootOrder) == 1
	})
	if len(bootDisks) > 0 {
		rootDiskName = bootDisks[0].Name
	}

	rootVols := lo.Filter(vm.Spec.Template.Spec.Volumes, func(volume kubevirtv1.Volume, idx int) bool {
		return volume.Name == rootDiskName
	})
	if len(rootVols) == 0 {
		if err := pkg.PatchAnnotations(ctx, reconciler.Client, &vm, annotation.VinkOperatingSystem.Name, struct{}{}); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to patch VirtualMachine operating system: %w", err)
		}
		return ctrl.Result{}, nil
	}

	rootDv := cdiv1beta1.DataVolume{}
	err := reconciler.Client.Get(ctx, types.NamespacedName{Namespace: vm.Namespace, Name: rootVols[0].DataVolume.Name}, &rootDv)
	if err != nil && !apierr.IsNotFound(err) {
		return ctrl.Result{}, fmt.Errorf("failed to get DataVolume: %w", err)
	}
	if apierr.IsNotFound(err) || rootDv.Annotations == nil {
		if err := pkg.PatchAnnotations(ctx, reconciler.Client, &vm, annotation.VinkOperatingSystem.Name, struct{}{}); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to patch VirtualMachine operating system: %w", err)
		}
		return ctrl.Result{}, nil
	}

	operatingSystem := apitypes.OperatingSystem{
		Name:    rootDv.Labels[label.VinkOperatingSystem.Name],
		Version: rootDv.Labels[label.VinkOperatingSystemVersion.Name],
	}

	if err := pkg.PatchAnnotations(ctx, reconciler.Client, &vm, annotation.VinkOperatingSystem.Name, &operatingSystem); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to patch VirtualMachine operating system: %w", err)
	}

	// operatingSystemBys, err := json.Marshal(&operatingSystem)
	// if err != nil {
	// 	return ctrl.Result{}, fmt.Errorf("marshal network info failed: %w", err)
	// }

	// patch := map[string]interface{}{
	// 	"metadata": map[string]interface{}{
	// 		"annotations": map[string]string{
	// 			annotation.VinkOperatingSystem.Name: string(operatingSystemBys),
	// 		},
	// 	},
	// }
	// patchBytes, err := json.Marshal(patch)
	// if err != nil {
	// 	return ctrl.Result{}, fmt.Errorf("failed to marshal patch: %w", err)
	// }

	// if err := reconciler.Client.Patch(ctx, &vm, client.RawPatch(types.MergePatchType, patchBytes)); err != nil {
	// 	return ctrl.Result{}, fmt.Errorf("failed to patch VirtualMachine operating system: %w", err)
	// }

	return ctrl.Result{}, nil
}

func (reconciler *OperatingSystemReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("virtualmachine_operating_system").
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
