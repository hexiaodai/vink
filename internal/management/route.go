package management

import (
	"github.com/gorilla/mux"
	resource_v1alpha1 "github.com/kubevm.io/vink/apis/management/resource/v1alpha1"
	vmv1alpha1 "github.com/kubevm.io/vink/apis/management/virtualmachine/v1alpha1"
	"github.com/kubevm.io/vink/internal/management/resource"
	"github.com/kubevm.io/vink/internal/management/virtualmachine"
	"github.com/kubevm.io/vink/pkg/informer"
	"google.golang.org/grpc/reflection"
)

func RegisterGRPCRoutes(kubeInformerFactory informer.KubeInformerFactory) (func(s reflection.GRPCServer), error) {
	return func(s reflection.GRPCServer) {
		resource_v1alpha1.RegisterResourceWatchManagementServer(s, resource.NewResourceWatchManagement(kubeInformerFactory))
		resource_v1alpha1.RegisterResourceManagementServer(s, resource.NewResourceManagement())
		vmv1alpha1.RegisterVirtualMachineManagementServer(s, virtualmachine.NewVirtualMachineManagement())
		reflection.Register(s)
	}, nil
}

func RegisterHTTPRoutes() (func(r *mux.Router), error) {
	return func(router *mux.Router) {
		virtualmachine.RegisterSerialConsole(router)
	}, nil
}
