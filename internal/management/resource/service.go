package resource

import (
	"fmt"

	resource_v1alpha1 "github.com/kubevm.io/vink/apis/management/resource/v1alpha1"
	"github.com/kubevm.io/vink/internal/management/resource/business"
	pkg_cache "github.com/kubevm.io/vink/internal/pkg/cache"
	"github.com/kubevm.io/vink/pkg/utils"
	"k8s.io/apimachinery/pkg/watch"
)

func NewResourceListWatchManagement(cache pkg_cache.Cache, subscribe pkg_cache.Subscribe) resource_v1alpha1.ResourceListWatchManagementServer {
	return &resourceListWatchManagement{cache: cache, subscribe: subscribe}
}

type resourceListWatchManagement struct {
	cache     pkg_cache.Cache
	subscribe pkg_cache.Subscribe

	resource_v1alpha1.UnsafeResourceListWatchManagementServer
}

func (rlw *resourceListWatchManagement) ListWatch(request *resource_v1alpha1.ListWatchRequest, server resource_v1alpha1.ResourceListWatchManagement_ListWatchServer) error {
	gvr := utils.ConvertGVR(request.GroupVersionResource)

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

	event := rlw.subscribe.Subscribe(gvr, metadatas)
	defer rlw.subscribe.Unsubscribe(gvr, event)

	for {
		select {
		case <-server.Context().Done():
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
				return err
			}
		}
	}
}
