package management

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	resource_v1alpha1 "github.com/kubevm.io/vink/apis/management/resource/v1alpha1"
	vmv1alpha1 "github.com/kubevm.io/vink/apis/management/virtualmachine/v1alpha1"
	"github.com/kubevm.io/vink/internal/management/resource"
	"github.com/kubevm.io/vink/internal/management/virtualmachine"
	"github.com/kubevm.io/vink/internal/pkg/cache"
	"github.com/kubevm.io/vink/pkg/clients"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RegisterHTTPRoutes() []func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error {
	return []func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error{}
}

func RegisterGRPCRoutes(cli clients.Clients, cache cache.Cache, subscribe cache.Subscribe) (func(s reflection.GRPCServer), error) {
	return func(s reflection.GRPCServer) {
		resource_v1alpha1.RegisterResourceListWatchManagementServer(s, resource.NewResourceListWatchManagement(cache, subscribe))
		vmv1alpha1.RegisterVirtualMachineManagementServer(s, virtualmachine.NewResourceListWatchManagement(cli))
		reflection.Register(s)
	}, nil
}
