package namespace

import (
	"context"

	"github.com/hexiaodai/vink/internal/management/namespace/business"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/api/errors"
	nsv1alpha1 "vink.io/api/management/namespace/v1alpha1"
)

func NewNamespaceManagement() nsv1alpha1.NamespaceManagementServer {
	return &dataVolumeManagement{}
}

type dataVolumeManagement struct {
	nsv1alpha1.UnimplementedNamespaceManagementServer
}

func (dvm *dataVolumeManagement) CreateNamespace(ctx context.Context, request *nsv1alpha1.CreateNamespaceRequest) (*nsv1alpha1.Namespace, error) {
	dv, err := business.CreateNamespace(ctx, request.Name, request.NamespaceConfig)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create namespace: %v", err)
	}
	return dv, nil
}

func (dvm *dataVolumeManagement) DeleteNamespace(ctx context.Context, request *nsv1alpha1.DeleteNamespaceRequest) (*nsv1alpha1.DeleteNamespaceResponse, error) {
	if err := business.DeleteNamespace(ctx, request.Name); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete namespace: %v", err)
	}
	return &nsv1alpha1.DeleteNamespaceResponse{}, nil
}

func (dvm *dataVolumeManagement) ListNamespaces(ctx context.Context, request *nsv1alpha1.ListNamespacesRequest) (*nsv1alpha1.ListNamespacesResponse, error) {
	namespaces, opts, err := business.ListNamespaces(ctx, request.Options)
	if errors.IsResourceExpired(err) {
		return nil, status.Errorf(codes.NotFound, "resource expired: %v", err)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list namespace: %v", err)
	}
	return &nsv1alpha1.ListNamespacesResponse{
		Items:   namespaces,
		Options: opts,
	}, nil
}
