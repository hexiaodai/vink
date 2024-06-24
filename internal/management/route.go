package management

import (
	"context"

	dvv1alpha1 "github.com/kubevm.io/vink/apis/management/datavolume/v1alpha1"
	nsv1alpha1 "github.com/kubevm.io/vink/apis/management/namespace/v1alpha1"
	nwv1alpha1 "github.com/kubevm.io/vink/apis/management/network/v1alpha1"
	vmv1alpha1 "github.com/kubevm.io/vink/apis/management/virtualmachine/v1alpha1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kubevm.io/vink/internal/management/datavolume"
	"github.com/kubevm.io/vink/internal/management/namespace"
	"github.com/kubevm.io/vink/internal/management/network"
	"github.com/kubevm.io/vink/internal/management/virtualmachine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RegisterHTTPRoutes() []func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error {
	return []func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error{
		dvv1alpha1.RegisterDataVolumeManagementHandler,
		nwv1alpha1.RegisterNetworkManagementHandler,
		vmv1alpha1.RegisterVirtualMachineManagementHandler,
		nsv1alpha1.RegisterNamespaceManagementHandler,
	}
}

func RegisterGRPCRoutes() (func(s reflection.GRPCServer), error) {
	return func(s reflection.GRPCServer) {
		dvv1alpha1.RegisterDataVolumeManagementServer(s, datavolume.NewDataVolumeManagement())
		nwv1alpha1.RegisterNetworkManagementServer(s, network.NewNetworkManagement())
		vmv1alpha1.RegisterVirtualMachineManagementServer(s, virtualmachine.NewVirtualMachineManagement())
		nsv1alpha1.RegisterNamespaceManagementServer(s, namespace.NewNamespaceManagement())
		reflection.Register(s)
	}, nil
}
