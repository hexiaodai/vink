package management

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	resource_v1alpha1 "github.com/kubevm.io/vink/apis/management/resource/v1alpha1"
	vmv1alpha1 "github.com/kubevm.io/vink/apis/management/virtualmachine/v1alpha1"
	"github.com/kubevm.io/vink/internal/management/resource"
	"github.com/kubevm.io/vink/internal/management/virtualmachine"
	resource_event_listener "github.com/kubevm.io/vink/internal/pkg/resource-event-listener"
	"github.com/kubevm.io/vink/pkg/informer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RegisterHTTPRoutes() []func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error {
	return []func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error{}
}

func RegisterGRPCRoutes(kubeInformerFactory informer.KubeInformerFactory, resourceEventListener resource_event_listener.ResourceEventListener) (func(s reflection.GRPCServer), error) {
	return func(s reflection.GRPCServer) {
		resource_v1alpha1.RegisterResourceListWatchManagementServer(s, resource.NewResourceListWatchManagement(kubeInformerFactory, resourceEventListener))
		resource_v1alpha1.RegisterResourceManagementServer(s, resource.NewResourceManagement())
		vmv1alpha1.RegisterVirtualMachineManagementServer(s, virtualmachine.NewVirtualMachineManagement())
		reflection.Register(s)
	}, nil
}
