package resource

import (
	"context"
	"fmt"

	"github.com/kubevm.io/vink/apis/apiextensions/v1alpha1"
	resource_v1alpha1 "github.com/kubevm.io/vink/apis/management/resource/v1alpha1"
	"github.com/kubevm.io/vink/internal/management/resource/business"
	resource_event_listener "github.com/kubevm.io/vink/internal/pkg/resource-event-listener"
	"github.com/kubevm.io/vink/pkg/clients/gvr"
	"github.com/kubevm.io/vink/pkg/informer"
	"github.com/kubevm.io/vink/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"k8s.io/apimachinery/pkg/watch"
)

func NewResourceListWatchManagement(kubeInformerFactory informer.KubeInformerFactory, resourceEventListener resource_event_listener.ResourceEventListener) resource_v1alpha1.ResourceListWatchManagementServer {
	return &resourceListWatchManagement{kubeInformerFactory: kubeInformerFactory, resourceEventListener: resourceEventListener}
}

type resourceListWatchManagement struct {
	resourceEventListener resource_event_listener.ResourceEventListener
	kubeInformerFactory   informer.KubeInformerFactory

	resource_v1alpha1.UnsafeResourceListWatchManagementServer
}

func (rlw *resourceListWatchManagement) ListWatch(request *resource_v1alpha1.ListWatchRequest, server resource_v1alpha1.ResourceListWatchManagement_ListWatchServer) error {
	gvr := gvr.ResolveGVR(request.GroupVersionResource)

	crds, metadatas, err := business.List(server.Context(), gvr, request.Options)
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
				crd, err := utils.ConvertToCustomResourceDefinition(e.Payload)
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

func NewResourceManagement() resource_v1alpha1.ResourceManagementServer {
	return &resourceManagement{}
}

type resourceManagement struct {
	resource_v1alpha1.UnsafeResourceManagementServer
}

func (r *resourceManagement) Create(ctx context.Context, request *resource_v1alpha1.CreateRequest) (*v1alpha1.CustomResourceDefinition, error) {
	gvr := gvr.ResolveGVR(request.GroupVersionResource)
	return business.Create(ctx, gvr, request.Data)
}

// Get implements v1alpha1.ResourceManagementServer.
func (r *resourceManagement) Get(context.Context, *resource_v1alpha1.GetRequest) (*v1alpha1.CustomResourceDefinition, error) {
	panic("unimplemented")
}

// Update implements v1alpha1.ResourceManagementServer.
func (r *resourceManagement) Update(context.Context, *resource_v1alpha1.UpdateRequest) (*v1alpha1.CustomResourceDefinition, error) {
	panic("unimplemented")
}

func (r *resourceManagement) Delete(ctx context.Context, request *resource_v1alpha1.DeleteRequest) (*emptypb.Empty, error) {
	gvr := gvr.ResolveGVR(request.GroupVersionResource)
	if err := business.Delete(ctx, gvr, request.NamespaceName); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
