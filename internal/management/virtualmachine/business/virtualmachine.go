package business

import (
	"context"

	vmv1alpha1 "github.com/kubevm.io/vink/apis/management/virtualmachine/v1alpha1"
	"github.com/kubevm.io/vink/apis/types"
	"github.com/kubevm.io/vink/pkg/clients"
	virtv1 "kubevirt.io/api/core/v1"
)

const (
	SerialConsoleRequestPathTmpl = "/apis/vink.io/v1alpha1/namespaces/{namespace}/virtualmachines/{name}/console"
)

func VirtualMachinePowerState(ctx context.Context, namespaceName *types.NamespaceName, powerState vmv1alpha1.VirtualMachinePowerStateRequest_PowerState) error {
	cli := clients.Instance.VirtualMachine(namespaceName.Namespace)

	switch powerState {
	case vmv1alpha1.VirtualMachinePowerStateRequest_ON:
		return cli.Start(ctx, namespaceName.Name, &virtv1.StartOptions{})
	case vmv1alpha1.VirtualMachinePowerStateRequest_OFF:
		return cli.Stop(ctx, namespaceName.Name, &virtv1.StopOptions{})
	case vmv1alpha1.VirtualMachinePowerStateRequest_REBOOT:
		return cli.Restart(ctx, namespaceName.Name, &virtv1.RestartOptions{})
	case vmv1alpha1.VirtualMachinePowerStateRequest_FORCE_OFF:
		return cli.ForceStop(ctx, namespaceName.Name, &virtv1.StopOptions{})
	case vmv1alpha1.VirtualMachinePowerStateRequest_FORCE_REBOOT:
		return cli.ForceRestart(ctx, namespaceName.Name, &virtv1.RestartOptions{})
	}

	return nil
}
