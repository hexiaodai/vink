package virtualmachine

import (
	"context"

	kubevirtv1 "kubevirt.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	LabelTpl = "vink.kubevm.io/tpl"

	LabelTplApplied = "vink.kubevm.io/tpl-applied"
)

type Reconciler struct {
	Client client.Client
	Cache  cache.Cache
}

func (r *Reconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("virtualmachine").
		For(&kubevirtv1.VirtualMachine{}).
		Complete(r)
}
