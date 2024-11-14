package resource

import (
	"context"
	"errors"
	"fmt"

	resource_v1alpha1 "github.com/kubevm.io/vink/apis/management/resource/v1alpha1"
	"github.com/kubevm.io/vink/internal/management/resource/business"
	"github.com/kubevm.io/vink/pkg/clients/gvr"
	"github.com/kubevm.io/vink/pkg/informer"
	"google.golang.org/protobuf/types/known/emptypb"
	"k8s.io/client-go/tools/cache"
)

func NewResourceWatchManagement(kubeInformerFactory informer.KubeInformerFactory) resource_v1alpha1.ResourceWatchManagementServer {
	return &resourceWatchManagement{kubeInformerFactory: kubeInformerFactory}
}

type resourceWatchManagement struct {
	kubeInformerFactory informer.KubeInformerFactory

	resource_v1alpha1.UnsafeResourceWatchManagementServer
}

func (rw *resourceWatchManagement) Watch(request *resource_v1alpha1.WatchRequest, server resource_v1alpha1.ResourceWatchManagement_WatchServer) error {
	gvr := gvr.ResolveGVR(request.ResourceType)

	informer, ok := rw.kubeInformerFactory.Informers()[gvr]
	if !ok {
		return fmt.Errorf("failed to find informer for %s", gvr.String())
	}

	filter, err := business.FilterFuncWithFieldSelector(request.Options)
	if err != nil {
		return fmt.Errorf("failed to filter field selector: %v", err)
	}
	filterFuncs := []business.FilterFunc{filter}

	errCh := make(chan error)

	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if err := business.SendResourceEvent(resource_v1alpha1.EventType_ADDED, obj, filterFuncs, server); err != nil {
				errCh <- err
			}
		},
		UpdateFunc: func(_, newObj interface{}) {
			if err := business.SendResourceEvent(resource_v1alpha1.EventType_MODIFIED, newObj, filterFuncs, server); err != nil {
				errCh <- err
			}
		},
		DeleteFunc: func(obj interface{}) {
			if err := business.SendResourceEvent(resource_v1alpha1.EventType_DELETED, obj, filterFuncs, server); err != nil {
				errCh <- err
			}
		},
	}
	registration, err := informer.AddEventHandler(eventHandler)
	if err != nil {
		return fmt.Errorf("failed to AddEventHandler, error: %v", err)
	}
	defer informer.RemoveEventHandler(registration)

	if err := server.Send(&resource_v1alpha1.WatchResponse{}); err != nil {
		return errors.New("failed to send readiness status to client")
	}

	select {
	case err := <-errCh:
		fmt.Println(err)
		return err
	case <-server.Context().Done():
		fmt.Println("stopping resource watch")
		return nil
	}
}

func NewResourceManagement() resource_v1alpha1.ResourceManagementServer {
	return &resourceManagement{}
}

type resourceManagement struct {
	resource_v1alpha1.UnsafeResourceManagementServer
}

func (r *resourceManagement) Create(ctx context.Context, request *resource_v1alpha1.CreateRequest) (*resource_v1alpha1.Resource, error) {
	gvr := gvr.ResolveGVR(request.ResourceType)
	crd, err := business.Create(ctx, gvr, request.Data)
	if err != nil {
		return nil, err
	}
	return &resource_v1alpha1.Resource{Data: crd}, nil
}

func (r *resourceManagement) Get(ctx context.Context, request *resource_v1alpha1.GetRequest) (*resource_v1alpha1.Resource, error) {
	gvr := gvr.ResolveGVR(request.ResourceType)
	crd, err := business.Get(ctx, gvr, request.NamespaceName.Namespace, request.NamespaceName.Name)
	if err != nil {
		return nil, err
	}
	return &resource_v1alpha1.Resource{Data: crd}, nil
}

func (r *resourceManagement) List(ctx context.Context, request *resource_v1alpha1.ListRequest) (*resource_v1alpha1.ListResponse, error) {
	gvr := gvr.ResolveGVR(request.ResourceType)
	crds, err := business.List(ctx, gvr, request.Options)
	if err != nil {
		return nil, err
	}
	return &resource_v1alpha1.ListResponse{Items: crds}, nil
}

func (r *resourceManagement) Update(ctx context.Context, request *resource_v1alpha1.UpdateRequest) (*resource_v1alpha1.Resource, error) {
	gvr := gvr.ResolveGVR(request.ResourceType)
	crd, err := business.Update(ctx, gvr, request.Data)
	if err != nil {
		return nil, err
	}
	return &resource_v1alpha1.Resource{Data: crd}, nil
}

func (r *resourceManagement) Delete(ctx context.Context, request *resource_v1alpha1.DeleteRequest) (*emptypb.Empty, error) {
	gvr := gvr.ResolveGVR(request.ResourceType)
	if err := business.Delete(ctx, gvr, request.NamespaceName); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
