package template_instance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kubevm.io/vink/pkg/k8s/apis/vink/v1alpha1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler struct {
	Client client.Client
	Cache  cache.Cache
}

func (r *Reconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	tplInstance := &v1alpha1.TemplateInstance{}
	if err := r.Client.Get(ctx, request.NamespacedName, tplInstance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if tplInstance.Status.Applied {
		return ctrl.Result{}, nil
	}

	var tpl v1alpha1.Template
	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: tplInstance.Namespace, Name: tplInstance.Spec.Template}, &tpl); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get Template %q in namespace %q for TemplateInstance %q: %w", tplInstance.Spec.Template, tplInstance.Namespace, tplInstance.Name, err)
	}

	vmCfg, err := r.buildVirtualMachineFromTemplate(ctx, &tpl, tplInstance)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to build VirtualMachine from Template %q: %w", tpl.Name, err)
	}

	data, _ := json.MarshalIndent(vmCfg, "", "  ")
	fmt.Println(string(data))

	if err := r.Client.Create(ctx, vmCfg); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create VirtualMachine %q: %w", vmCfg.Name, err)
	}

	tplInstance.Status.Applied = true
	if err := r.Client.Status().Update(ctx, tplInstance); err != nil {
		if apierr.IsConflict(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to update TemplateInstance %q: %w", tplInstance.Name, err)
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("templateinstance").
		For(&v1alpha1.TemplateInstance{}).
		Complete(r)
}
