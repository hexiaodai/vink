package virtualmachinesummarys

import (
	"context"
	"fmt"

	kubeovn "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubevm.io/vink/pkg/k8s/apis/vink/v1alpha1"
	"github.com/kubevm.io/vink/pkg/log"
	"github.com/samber/lo"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NetworkReconciler struct {
	Client client.Client
	Cache  cache.Cache
}

func (reconciler *NetworkReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	list := &v1alpha1.VirtualMachineSummaryList{}
	if err := reconciler.Client.List(ctx, list); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to list VirtualMachineSummarys: %w", err)
	}

	var ip kubeovn.IP
	if err := reconciler.Client.Get(ctx, request.NamespacedName, &ip); err != nil && !apierr.IsNotFound(err) {
		return ctrl.Result{}, fmt.Errorf("failed to get IP: %w", err)
	}

	if ip.Spec.PodType != "VirtualMachine" {
		return ctrl.Result{}, nil
	}

	for _, item := range list.Items {
		if ip.Spec.PodName != item.Name || ip.Spec.Namespace != item.Namespace {
			continue
		}

		newIPs := make(map[string]*v1alpha1.IP, len(item.Status.Network.IPs))
		for _, vmIP := range item.Status.Network.IPs {
			eq := (vmIP.ObjectMeta.Namespace == request.Namespace && vmIP.ObjectMeta.Name == request.Name)
			if eq && len(ip.UID) == 0 {
				continue
			}
			newIP := vmIP
			if eq && len(ip.UID) > 0 && vmIP.ObjectMeta.ResourceVersion != ip.ResourceVersion {
				newIP = v1alpha1.IPFromKubeOVN(&ip)
			}
			newIPs[string(newIP.ObjectMeta.UID)] = newIP
		}
		if _, ok := newIPs[string(ip.UID)]; !ok && len(string(ip.UID)) > 0 {
			newIPs[string(ip.UID)] = v1alpha1.IPFromKubeOVN(&ip)
		}

		item.Status.Network.IPs = lo.MapToSlice(newIPs, func(key string, value *v1alpha1.IP) *v1alpha1.IP {
			return value
		})

		if err := reconciler.Client.Status().Update(ctx, &item); err != nil {
			if apierr.IsConflict(err) {
				log.Debugf("VirtualMachineSummary %s/%s conflict", item.Namespace, item.Name)
				return ctrl.Result{Requeue: true}, nil
			}
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (reconciler *NetworkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeovn.IP{}).
		Complete(reconciler)
}
