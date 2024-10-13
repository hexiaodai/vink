package virtualmachinesummarys

import (
	"context"
	"fmt"

	"github.com/kubevm.io/vink/pkg/k8s/apis/vink/v1alpha1"
	"github.com/kubevm.io/vink/pkg/log"
	corev1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type HostReconciler struct {
	Client client.Client
	Cache  cache.Cache
}

func (reconciler *HostReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	var summaryList v1alpha1.VirtualMachineSummaryList
	err := reconciler.Cache.List(ctx, &summaryList)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to list VirtualMachineSummarys: %w", err)
	}

	var node corev1.Node
	if err := reconciler.Cache.Get(ctx, request.NamespacedName, &node); err != nil && !apierr.IsNotFound(err) {
		return ctrl.Result{}, fmt.Errorf("failed to get Node: %w", err)
	}

	for _, item := range summaryList.Items {
		if len(string(node.UID)) > 0 && item.Status.Host != nil && item.ResourceVersion == item.Status.Host.ObjectMeta.ResourceVersion {
			continue
		}

		if len(string(node.UID)) == 0 {
			item.Status.Host = nil
		} else {
			item.Status.Host = v1alpha1.NodeFromKube(&node)
		}

		if err := reconciler.Client.Status().Update(ctx, &item); err != nil {
			if apierr.IsConflict(err) {
				log.Debugf("VirtualMachineSummary %s/%s conflict", request.Namespace, request.Name)
				return ctrl.Result{Requeue: true}, nil
			}
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (reconciler *HostReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Node{}).
		Complete(reconciler)
}
