package management

import (
	"github.com/gorilla/mux"
	resource_v1alpha1 "github.com/kubevm.io/vink/apis/management/resource/v1alpha1"
	vmv1alpha1 "github.com/kubevm.io/vink/apis/management/virtualmachine/v1alpha1"
	"github.com/kubevm.io/vink/internal/management/resource"
	"github.com/kubevm.io/vink/internal/management/virtualmachine"
	resource_event_listener "github.com/kubevm.io/vink/internal/pkg/resource-event-listener"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/informer"
	"google.golang.org/grpc/reflection"
)

func RegisterGRPCRoutes(clients clients.Clients, kubeInformerFactory informer.KubeInformerFactory, resourceEventListener resource_event_listener.ResourceEventListener) (func(s reflection.GRPCServer), error) {
	return func(s reflection.GRPCServer) {
		resource_v1alpha1.RegisterResourceListWatchManagementServer(s, resource.NewResourceListWatchManagement(clients, kubeInformerFactory, resourceEventListener))
		resource_v1alpha1.RegisterResourceManagementServer(s, resource.NewResourceManagement(clients))
		vmv1alpha1.RegisterVirtualMachineManagementServer(s, virtualmachine.NewVirtualMachineManagement(clients))
		reflection.Register(s)
	}, nil
}

func RegisterHTTPRoutes(clients clients.Clients) (func(r *mux.Router), error) {
	return func(router *mux.Router) {
		virtualmachine.RegisterSerialConsole(router, clients)
	}, nil
}
