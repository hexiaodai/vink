package resource

import (
	"context"
	"fmt"

	resource_v1alpha1 "github.com/kubevm.io/vink/apis/management/resource/v1alpha1"
	"github.com/kubevm.io/vink/internal/management/resource/business"
	resource_event_listener "github.com/kubevm.io/vink/internal/pkg/resource-event-listener"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/clients/gvr"
	"github.com/kubevm.io/vink/pkg/informer"
	"github.com/kubevm.io/vink/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"k8s.io/apimachinery/pkg/watch"
)

func NewResourceListWatchManagement(clients clients.Clients, kubeInformerFactory informer.KubeInformerFactory, resourceEventListener resource_event_listener.ResourceEventListener) resource_v1alpha1.ResourceListWatchManagementServer {
	return &resourceListWatchManagement{clients: clients, kubeInformerFactory: kubeInformerFactory, resourceEventListener: resourceEventListener}
}

type resourceListWatchManagement struct {
	clients               clients.Clients
	resourceEventListener resource_event_listener.ResourceEventListener
	kubeInformerFactory   informer.KubeInformerFactory

	resource_v1alpha1.UnsafeResourceListWatchManagementServer
}

func (rlw *resourceListWatchManagement) ListWatch(request *resource_v1alpha1.ListWatchRequest, server resource_v1alpha1.ResourceListWatchManagement_ListWatchServer) error {
	gvr := gvr.ResolveGVR(request.ResourceType)

	crds, metadatas, err := business.List(server.Context(), rlw.clients, gvr, request.Options)
	if err != nil {
		return err
	}

	if err := server.Send(&resource_v1alpha1.ListWatchResponse{
		EventType: resource_v1alpha1.EventType_ADDED,
		Items:     crds,
	}); err != nil {
		return err
	}

	if !request.Options.Watch {
		return nil
	}

	event := rlw.resourceEventListener.AddSubscription(gvr, metadatas)
	defer rlw.resourceEventListener.RemoveSubscription(gvr, event)

	for {
		select {
		case <-server.Context().Done():
			fmt.Println("resource list watch canceled")
			return nil
		case e, ok := <-event:
			if !ok {
				return nil
			}

			resp := resource_v1alpha1.ListWatchResponse{}

			switch e.Type {
			case watch.Modified:
				crd, err := clients.InterfaceToJSON(e.Payload)
				if err != nil {
					return fmt.Errorf("failed to convert payload to CustomResourceDefinition: %v", err)
				}
				resp.EventType = resource_v1alpha1.EventType_MODIFIED
				resp.Items = append(resp.Items, crd)
			case watch.Deleted:
				nn, err := utils.ConvertToNamespaceName(e.Payload)
				if err != nil {
					return fmt.Errorf("failed to convert payload to NamespaceName: %v", err)
				}
				resp.EventType = resource_v1alpha1.EventType_DELETED
				resp.Deleted = nn
			}

			if err := server.Send(&resp); err != nil {
				fmt.Println("failed to send response to client")
				return err
			}
		}
	}
}

func NewResourceManagement(clients clients.Clients) resource_v1alpha1.ResourceManagementServer {
	return &resourceManagement{
		clients: clients,
	}
}

type resourceManagement struct {
	clients clients.Clients

	resource_v1alpha1.UnsafeResourceManagementServer
}

func (r *resourceManagement) Create(ctx context.Context, request *resource_v1alpha1.CreateRequest) (*resource_v1alpha1.CustomResourceDefinitionResponse, error) {
	gvr := gvr.ResolveGVR(request.ResourceType)
	crd, err := business.Create(ctx, r.clients, gvr, request.Data)
	if err != nil {
		return nil, err
	}
	return &resource_v1alpha1.CustomResourceDefinitionResponse{Data: crd}, nil
}

func (r *resourceManagement) Get(ctx context.Context, request *resource_v1alpha1.GetRequest) (*resource_v1alpha1.CustomResourceDefinitionResponse, error) {
	gvr := gvr.ResolveGVR(request.ResourceType)
	crd, err := business.Get(ctx, r.clients, gvr, request.NamespaceName.Namespace, request.NamespaceName.Name)
	if err != nil {
		return nil, err
	}
	return &resource_v1alpha1.CustomResourceDefinitionResponse{Data: crd}, nil
}

func (r *resourceManagement) Update(ctx context.Context, request *resource_v1alpha1.UpdateRequest) (*resource_v1alpha1.CustomResourceDefinitionResponse, error) {
	gvr := gvr.ResolveGVR(request.ResourceType)
	crd, err := business.Update(ctx, r.clients, gvr, request.Data)
	if err != nil {
		return nil, err
	}
	return &resource_v1alpha1.CustomResourceDefinitionResponse{Data: crd}, nil
}

func (r *resourceManagement) Delete(ctx context.Context, request *resource_v1alpha1.DeleteRequest) (*emptypb.Empty, error) {
	gvr := gvr.ResolveGVR(request.ResourceType)
	if err := business.Delete(ctx, r.clients, gvr, request.NamespaceName); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
