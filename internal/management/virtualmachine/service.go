package virtualmachine

import (
	"context"

	vmv1alpha1 "github.com/kubevm.io/vink/apis/management/virtualmachine/v1alpha1"
	"github.com/kubevm.io/vink/internal/management/virtualmachine/business"
	"github.com/kubevm.io/vink/pkg/clients"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewResourceListWatchManagement(clients clients.Clients) vmv1alpha1.VirtualMachineManagementServer {
	return &virtualMachineManagement{
		clients: clients,
	}
}

type virtualMachineManagement struct {
	clients clients.Clients

	vmv1alpha1.UnimplementedVirtualMachineManagementServer
}

func (m *virtualMachineManagement) VirtualMachinePowerState(ctx context.Context, request *vmv1alpha1.VirtualMachinePowerStateRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, business.VirtualMachinePowerState(ctx, m.clients, request.NamespaceName, request.PowerState)
}
