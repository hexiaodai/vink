package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kubevm.io/vink/apis/annotation"
	"github.com/kubevm.io/vink/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	kubevirtv1 "kubevirt.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type VMIHostReconciler struct {
	Client client.Client
	Cache  cache.Cache
}

func (reconciler *VMIHostReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	var vmi kubevirtv1.VirtualMachineInstance
	if err := reconciler.Client.Get(ctx, request.NamespacedName, &vmi); err != nil {
		if apierr.IsNotFound(err) {
			return ctrl.Result{}, nil
		} else {
			return ctrl.Result{}, fmt.Errorf("failed to get VirtualMachine: %w", err)
		}
	}

	if len(vmi.Status.NodeName) == 0 || vmi.DeletionTimestamp != nil {
		return ctrl.Result{}, nil
	}

	node := corev1.Node{}
	if err := reconciler.Client.Get(ctx, client.ObjectKey{Name: vmi.Status.NodeName}, &node); err != nil {
		return ctrl.Result{}, nil
	}

	newIPs := make([]string, 0)
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP || addr.Type == corev1.NodeExternalIP && len(addr.Address) > 0 {
			newIPs = append(newIPs, addr.Address)
		}
	}

	if vmi.Annotations == nil {
		vmi.Annotations = make(map[string]string)
	}

	oldIPs := make([]string, 0)
	_ = json.Unmarshal([]byte(vmi.Annotations[annotation.VinkHost.Name]), &oldIPs)

	if utils.CompareArrays(newIPs, oldIPs) {
		return ctrl.Result{}, nil
	}

	value, err := json.Marshal(newIPs)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("marshal error: %w", err)
	}

	vmi.Annotations[annotation.VinkHost.Name] = string(value)
	if err := reconciler.Client.Update(ctx, &vmi); err != nil {
		if apierr.IsConflict(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (reconciler *VMIHostReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("vmi_host").
		For(&kubevirtv1.VirtualMachineInstance{}).
		Complete(reconciler)
}
