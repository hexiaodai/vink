package virtualmachinesummarys

import (
	"context"
	"fmt"

	"github.com/kubevm.io/vink/pkg/k8s/apis/vink/v1alpha1"
	"github.com/kubevm.io/vink/pkg/log"
	"github.com/samber/lo"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DataVolumeReconciler struct {
	Client client.Client
	Cache  cache.Cache
}

func (reconciler *DataVolumeReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	var summaryList v1alpha1.VirtualMachineSummaryList
	if err := reconciler.Cache.List(ctx, &summaryList, &client.ListOptions{Namespace: request.Namespace}); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to list VirtualMachineSummarys: %w", err)
	}

	var dv cdiv1beta1.DataVolume
	if err := reconciler.Cache.Get(ctx, request.NamespacedName, &dv); err != nil && !apierr.IsNotFound(err) {
		return ctrl.Result{}, fmt.Errorf("failed to get DataVolume: %w", err)
	}

	for _, item := range summaryList.Items {
		changed := false
		newVmdvs := make(map[string]*v1alpha1.DataVolume, len(item.Status.DataVolumes))
		for _, vmdv := range item.Status.DataVolumes {
			eq := (vmdv.ObjectMeta.Namespace == request.Namespace && vmdv.ObjectMeta.Name == request.Name)
			if eq && len(string(dv.UID)) == 0 {
				changed = true
				continue
			}
			newVmdv := vmdv
			if eq && len(dv.UID) > 0 && vmdv.ObjectMeta.ResourceVersion != dv.ResourceVersion {
				changed = true
				newVmdv = v1alpha1.DataVolumeFromKubeVirt(&dv)
			}
			newVmdvs[string(newVmdv.ObjectMeta.UID)] = newVmdv
		}

		found := false
		for _, vol := range item.Status.VirtualMachine.Spec.Template.Spec.Volumes {
			if vol.DataVolume == nil {
				continue
			}
			if vol.Name == dv.Name {
				found = true
				break
			}
		}

		if _, ok := newVmdvs[string(dv.UID)]; !ok && found && len(string(dv.UID)) > 0 {
			changed = true
			newVmdvs[string(dv.UID)] = v1alpha1.DataVolumeFromKubeVirt(&dv)
		}
		if !changed {
			continue
		}

		item.Status.DataVolumes = lo.MapToSlice(newVmdvs, func(key string, value *v1alpha1.DataVolume) *v1alpha1.DataVolume {
			return value
		})
		if err := reconciler.Client.Status().Update(ctx, &item); err != nil {
			if apierr.IsConflict(err) {
				log.Debugf("VirtualMachineSummary %s/%s conflict", item.Namespace, item.Name)
				return ctrl.Result{Requeue: true}, nil
			}
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (reconciler *DataVolumeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cdiv1beta1.DataVolume{}).
		Complete(reconciler)
}
