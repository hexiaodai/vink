package virtualmachine

import (
	"context"

	vmv1alpha1 "github.com/kubevm.io/vink/apis/management/virtualmachine/v1alpha1"
	"github.com/kubevm.io/vink/internal/management/virtualmachine/business"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/api/errors"
)

func NewVirtualMachineManagement() vmv1alpha1.VirtualMachineManagementServer {
	return &virtualMachineManagement{}
}

type virtualMachineManagement struct {
	vmv1alpha1.UnimplementedVirtualMachineManagementServer
}

func (vmm *virtualMachineManagement) CreateVirtualMachine(ctx context.Context, request *vmv1alpha1.CreateVirtualMachineRequest) (*vmv1alpha1.VirtualMachine, error) {
	vm, err := business.CreateVirtualMachine(ctx, request.Namespace, request.Name, request.Config)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create virtual machine: %v", err)
	}
	return vm, nil
}

func (vmm *virtualMachineManagement) DeleteVirtualMachine(ctx context.Context, request *vmv1alpha1.DeleteVirtualMachineRequest) (*vmv1alpha1.DeleteVirtualMachineResponse, error) {
	if err := business.DeleteVirtualMachine(ctx, request.Namespace, request.Name); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete virtual machine: %v", err)
	}
	return &vmv1alpha1.DeleteVirtualMachineResponse{}, nil
}

func (vmm *virtualMachineManagement) ListVirtualMachines(ctx context.Context, request *vmv1alpha1.ListVirtualMachinesRequest) (*vmv1alpha1.ListVirtualMachinesResponse, error) {
	vms, opts, err := business.ListVirtualMachines(ctx, request.Namespace, request.Options)
	if errors.IsResourceExpired(err) {
		return nil, status.Errorf(codes.NotFound, "resource expired: %v", err)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list virtual machines: %v", err)
	}
	return &vmv1alpha1.ListVirtualMachinesResponse{
		Items:   vms,
		Options: opts,
	}, nil
}

func (vmm *virtualMachineManagement) ManageVirtualMachinePowerState(ctx context.Context, request *vmv1alpha1.ManageVirtualMachinePowerStateRequest) (*vmv1alpha1.VirtualMachine, error) {
	vm, err := business.ManageVirtualMachinePowerState(ctx, request.Namespace, request.Name, request.PowerState)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to manage virtual machine power state: %v", err)
	}
	return vm, nil
}
