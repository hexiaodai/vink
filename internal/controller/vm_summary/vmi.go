package summarys

import (
	"context"
	"fmt"

	"github.com/kubevm.io/vink/pkg/k8s/apis/vink/v1alpha1"
	"github.com/kubevm.io/vink/pkg/log"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	kubevirtv1 "kubevirt.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type VirtualMachineInstanceReconciler struct {
	Client client.Client
	Cache  cache.Cache
}

func (reconciler *VirtualMachineInstanceReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	var summary v1alpha1.VirtualMachineSummary
	err := reconciler.Cache.Get(ctx, request.NamespacedName, &summary)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get VirtualMachineSummary: %w", err)
	}

	var vmi kubevirtv1.VirtualMachineInstance
	if err := reconciler.Cache.Get(ctx, request.NamespacedName, &vmi); err != nil && !apierr.IsNotFound(err) {
		return ctrl.Result{}, fmt.Errorf("failed to get VirtualMachineInstance: %w", err)
	}
	if len(string(vmi.UID)) > 0 && summary.Status.VirtualMachineInstance != nil && vmi.ResourceVersion == summary.Status.VirtualMachineInstance.ObjectMeta.ResourceVersion {
		return ctrl.Result{}, nil
	}

	if len(string(vmi.UID)) == 0 {
		summary.Status.VirtualMachineInstance = nil
	} else {
		summary.Status.VirtualMachineInstance = v1alpha1.VirtualMachineInstanceFromKubeVirt(&vmi)
	}

	if err := reconciler.Client.Status().Update(ctx, &summary); err != nil {
		if apierr.IsConflict(err) {
			log.Debugf("VirtualMachineSummary %s/%s conflict", request.Namespace, request.Name)
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (reconciler *VirtualMachineInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("summarys_virtualmachineinstance").
		For(&kubevirtv1.VirtualMachineInstance{}).
		Complete(reconciler)
}
