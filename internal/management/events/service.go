package events

import (
	"context"

	events_v1alpha1 "github.com/kubevm.io/vink/apis/management/events/v1alpha1"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/informer"
)

func NewEventsManagement(clients clients.Clients, kubeInformerFactory informer.KubeInformerFactory) events_v1alpha1.EventsManagementServer {
	return &eventsManagement{clients: clients, kubeInformerFactory: kubeInformerFactory}
}

type eventsManagement struct {
	clients             clients.Clients
	kubeInformerFactory informer.KubeInformerFactory

	events_v1alpha1.UnsafeEventsManagementServer
}

func (m *eventsManagement) List(ctx context.Context, request *events_v1alpha1.ListEventsRequest) (*events_v1alpha1.Events, error) {
	return nil, nil
}
