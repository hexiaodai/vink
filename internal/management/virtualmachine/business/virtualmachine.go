package business

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubevm.io/vink/apis/common"
	vmv1alpha1 "github.com/kubevm.io/vink/apis/management/virtualmachine/v1alpha1"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/clients/gvr"
	"github.com/kubevm.io/vink/pkg/utils"
	"github.com/samber/lo"
	virtv1 "kubevirt.io/api/core/v1"
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
	if err := setupVirtualMachineRootVolume(vm, config.Storage.Root); err != nil {
		return nil, err
	}
	if err := setupVirtualMachineNetwork(vm, config.Network); err != nil {
		return nil, err
	}
	if err := setupVirtualMachineDataVolumes(vm, config.Storage.Data); err != nil {
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
