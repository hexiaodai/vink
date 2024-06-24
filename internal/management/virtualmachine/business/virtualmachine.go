package business

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hexiaodai/vink/pkg/clients"
	"github.com/hexiaodai/vink/pkg/clients/gvr"
	"github.com/hexiaodai/vink/pkg/utils"
	"github.com/samber/lo"
	virtv1 "kubevirt.io/api/core/v1"
	"vink.io/api/common"
	vmv1alpha1 "vink.io/api/management/virtualmachine/v1alpha1"
)

func ListVirtualMachines(ctx context.Context, namespace string, opts *common.ListOptions) ([]*vmv1alpha1.VirtualMachine, *common.ListOptions, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()

	unobj, err := dcli.Resource(gvr.From(virtv1.VirtualMachine{})).Namespace(namespace).List(ctx, utils.ConvertToK8sListOptions(opts))
	if err != nil {
		return nil, nil, err
	}
	obj, err := clients.FromUnstructuredList[virtv1.VirtualMachineList](unobj)
	if err != nil {
		return nil, nil, err
	}

	vms := make([]*vmv1alpha1.VirtualMachine, 0, len(obj.Items))
	for _, vm := range obj.Items {
		apivm, err := crdToAPIVirtualMachine(ctx, &vm)
		if err != nil {
			return nil, nil, err
		}
		vms = append(vms, apivm)
	}

	return vms, utils.ConvertToAPIListOptions(opts, obj.ListMeta), nil
}

func CreateVirtualMachine(ctx context.Context, namespace, name string, config *vmv1alpha1.VirtualMachineConfig) (*vmv1alpha1.VirtualMachine, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()

	vm := newSampleVirtualMachine(namespace, name, config)
	if err := setupVirtualMachineBootDisk(vm, config.Storage.BootDisk); err != nil {
		return nil, err
	}
	if err := setupVirtualMachineNetwork(vm, config.Network); err != nil {
		return nil, err
	}
	if err := setupVirtualMachineDataDisks(vm, config.Storage.DataDisks); err != nil {
		return nil, err
	}
	if err := setupVirtualMachineUserConfig(vm, config.UserConfig); err != nil {
		return nil, err
	}

	un, _ := clients.Unstructured(vm)
	unObj, err := dcli.Resource(gvr.From(virtv1.VirtualMachine{})).Namespace(namespace).Create(ctx, un, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	obj, err := clients.FromUnstructured[virtv1.VirtualMachine](unObj)
	if err != nil {
		return nil, err
	}

	return crdToAPIVirtualMachine(ctx, obj)
}

func DeleteVirtualMachine(ctx context.Context, namespace, name string) error {
	dcli := clients.GetClients().GetDynamicKubeClient()

	err := dcli.Resource(gvr.From(virtv1.VirtualMachine{})).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func ManageVirtualMachinePowerState(ctx context.Context, namespace, name string, powerState vmv1alpha1.ManageVirtualMachinePowerStateRequest_PowerState) (*vmv1alpha1.VirtualMachine, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()

	unobj, err := dcli.Resource(gvr.From(virtv1.VirtualMachine{})).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	obj, err := clients.FromUnstructured[virtv1.VirtualMachine](unobj)
	if err != nil {
		return nil, err
	}
	apivm, err := crdToAPIVirtualMachine(ctx, obj)
	if err != nil {
		return nil, err
	}

	if powerState == vmv1alpha1.ManageVirtualMachinePowerStateRequest_UNSPECIFIED {
		return apivm, nil
	}

	power := virtv1.RunStrategyUnknown
	switch powerState {
	case vmv1alpha1.ManageVirtualMachinePowerStateRequest_ON:
		power = virtv1.RunStrategyAlways
	case vmv1alpha1.ManageVirtualMachinePowerStateRequest_OFF:
		power = virtv1.RunStrategyHalted
	}

	if lo.FromPtr(obj.Spec.RunStrategy) == power {
		return apivm, nil
	}

	obj.Spec.RunStrategy = &power
	un, _ := clients.Unstructured(obj)
	unObj, err := dcli.Resource(gvr.From(virtv1.VirtualMachine{})).Namespace(namespace).Update(ctx, un, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	obj, err = clients.FromUnstructured[virtv1.VirtualMachine](unObj)
	if err != nil {
		return nil, err
	}

	apivm, err = crdToAPIVirtualMachine(ctx, obj)
	if err != nil {
		return nil, err
	}

	return apivm, nil
}

// func DeleteVirtualMachine(ctx context.Context, namespace, name string) error {
// 	dcli := clients.GetClients().GetDynamicKubeClient()

// 	unObj, err := dcli.Resource(gvr.From(virtv1.VirtualMachine{})).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
// 	if err != nil {
// 		return err
// 	}
// 	obj, err := clients.FromUnstructured[virtv1.VirtualMachine](unObj)
// 	if err != nil {
// 		return err
// 	}

// 	datavolumes := make([]*cdiv1beta1.DataVolume, 0, len(obj.Spec.Template.Spec.Volumes))
// 	exclude := map[string]struct{}{"cloudinit": {}, "boot": {}}
// 	for _, volume := range obj.Spec.Template.Spec.Volumes {
// 		if _, ok := exclude[volume.Name]; ok {
// 			continue
// 		}
// 		unpvcObj, err := dcli.Resource(gvr.From(corev1.PersistentVolumeClaim{})).Namespace(namespace).Get(ctx, volume.PersistentVolumeClaim.ClaimName, metav1.GetOptions{})
// 		if err != nil && errors.IsNotFound(err) {
// 			continue
// 		} else if err != nil {
// 			return err
// 		}
// 		pvcObj, err := clients.FromUnstructured[corev1.PersistentVolumeClaim](unpvcObj)
// 		if err != nil {
// 			return err
// 		}
// 		var ownerObj *metav1.OwnerReference
// 		for _, owner := range pvcObj.ObjectMeta.OwnerReferences {
// 			if owner.Kind == "DataVolume" {
// 				ownerObj = &owner
// 				break
// 			}
// 		}
// 		if ownerObj == nil {
// 			continue
// 		}
// 		undvObj, err := dcli.Resource(gvr.From(cdiv1beta1.DataVolume{})).Namespace(namespace).Get(ctx, ownerObj.Name, metav1.GetOptions{})
// 		if err != nil && errors.IsNotFound(err) {
// 			continue
// 		} else if err != nil {
// 			return err
// 		}
// 		dvObj, err := clients.FromUnstructured[cdiv1beta1.DataVolume](undvObj)
// 		if err != nil {
// 			return err
// 		}
// 		datavolumes = append(datavolumes, dvObj)
// 	}

// 	for _, dv := range datavolumes {
// 		newOwners := lo.Filter(dv.OwnerReferences, func(item metav1.OwnerReference, _ int) bool {
// 			return item.Kind != "VirtualMachine"
// 		})
// 		dv.OwnerReferences = newOwners
// 		undv, _ := clients.Unstructured(dv)
// 		_, err = dcli.Resource(gvr.From(cdiv1beta1.DataVolume{})).Namespace(namespace).Update(ctx, undv, metav1.UpdateOptions{})
// 		if err != nil && errors.IsNotFound(err) {
// 			continue
// 		} else if err != nil {
// 			return err
// 		}
// 	}

// 	err = dcli.Resource(gvr.From(virtv1.VirtualMachine{})).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
// 	if errors.IsNotFound(err) {
// 		return nil
// 	}
// 	return err
// }
