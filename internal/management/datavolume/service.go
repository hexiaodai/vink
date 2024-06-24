package datavolume

import (
	"context"

	"github.com/hexiaodai/vink/internal/management/datavolume/business"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/api/errors"
	dvv1alpha1 "vink.io/api/management/datavolume/v1alpha1"
)

func NewDataVolumeManagement() dvv1alpha1.DataVolumeManagementServer {
	return &dataVolumeManagement{}
}

type dataVolumeManagement struct {
	dvv1alpha1.UnimplementedDataVolumeManagementServer
}

func (dvm *dataVolumeManagement) CreateDataVolume(ctx context.Context, request *dvv1alpha1.CreateDataVolumeRequest) (*dvv1alpha1.DataVolume, error) {
	dv, err := business.CreateDataVolumes(ctx, request.Namespace, request.Name, request.Config)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create data volume: %v", err)
	}
	return dv, nil
}

func (dvm *dataVolumeManagement) DeleteDataVolume(ctx context.Context, request *dvv1alpha1.DeleteDataVolumeRequest) (*dvv1alpha1.DeleteDataVolumeResponse, error) {
	if err := business.DeleteDataVolume(ctx, request.Namespace, request.Name); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete data volume: %v", err)
	}
	return &dvv1alpha1.DeleteDataVolumeResponse{}, nil
}

func (dvm *dataVolumeManagement) ListDataVolumes(ctx context.Context, request *dvv1alpha1.ListDataVolumesRequest) (*dvv1alpha1.ListDataVolumesResponse, error) {
	datavolumes, opts, err := business.ListDataVolumes(ctx, request.Namespace, request.Options)
	if errors.IsResourceExpired(err) {
		// TODO: handle expired error
		// 前端应该清空 Continue 值，然后重试
		return nil, status.Errorf(codes.NotFound, "resource expired: %v", err)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list data volumes: %v", err)
	}
	return &dvv1alpha1.ListDataVolumesResponse{
		Items:   datavolumes,
		Options: opts,
	}, nil
}
