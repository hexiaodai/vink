package network

import (
	"context"

	"github.com/hexiaodai/vink/internal/management/network/business"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/api/errors"
	nwv1alpha1 "vink.io/api/management/network/v1alpha1"
)

func NewNetworkManagement() nwv1alpha1.NetworkManagementServer {
	return &networkManagement{}
}

type networkManagement struct {
	nwv1alpha1.UnimplementedNetworkManagementServer
}

func (nm *networkManagement) ListNodesNetworkInterfaces(ctx context.Context, request *nwv1alpha1.ListNodesNetworkInterfacesRequest) (*nwv1alpha1.ListNodesNetworkInterfacesResponse, error) {
	nfs, opts, err := business.ListNodesNetworkInterfaces(ctx, request.Options)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list node network interfaces: %v", err)
	}
	return &nwv1alpha1.ListNodesNetworkInterfacesResponse{
		Items:   nfs,
		Options: opts,
	}, nil
}

func (nm *networkManagement) CreateMultusConfig(ctx context.Context, request *nwv1alpha1.CreateMultusConfigRequest) (*nwv1alpha1.MultusConfig, error) {
	multusCfg, err := business.CreateMultusConfig(ctx, request.Name, request.Nic)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create multus config: %v", err)
	}
	return multusCfg, nil
}

func (nm *networkManagement) CreateSubnet(ctx context.Context, request *nwv1alpha1.CreateSubnetRequest) (*nwv1alpha1.Subnet, error) {
	subnet, err := business.CreateSubnet(ctx, request.Name, request.Config)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create subnet: %v", err)
	}
	return subnet, nil
}

func (nm *networkManagement) CreateIPPool(ctx context.Context, request *nwv1alpha1.CreateIPPoolRequest) (*nwv1alpha1.IPPool, error) {
	ippool, err := business.CreateIPPool(ctx, request.Name, request.Config)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create ippool: %v", err)
	}
	return ippool, nil
}

func (nm *networkManagement) UpdateMultusConfig(ctx context.Context, request *nwv1alpha1.UpdateMultusConfigRequest) (*nwv1alpha1.MultusConfig, error) {
	return nil, nil
}

func (nm *networkManagement) UpdateSubnet(ctx context.Context, request *nwv1alpha1.UpdateSubnetRequest) (*nwv1alpha1.Subnet, error) {
	return nil, nil
}

func (nm *networkManagement) UpdateIPPool(ctx context.Context, request *nwv1alpha1.UpdateIPPoolRequest) (*nwv1alpha1.IPPool, error) {
	return nil, nil
}

func (nm *networkManagement) DeleteMultusConfig(ctx context.Context, request *nwv1alpha1.DeleteMultusConfigRequest) (*nwv1alpha1.DeleteMultusConfigResponse, error) {
	if err := business.DeleteMultusConfig(ctx, request.Name); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete multus config: %v", err)
	}
	return &nwv1alpha1.DeleteMultusConfigResponse{}, nil
}

func (nm *networkManagement) DeleteSubnet(ctx context.Context, request *nwv1alpha1.DeleteSubnetRequest) (*nwv1alpha1.DeleteSubnetResponse, error) {
	if err := business.DeleteSubnet(ctx, request.Name); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete subnet: %v", err)
	}
	return &nwv1alpha1.DeleteSubnetResponse{}, nil
}

func (nm *networkManagement) DeleteIPPool(ctx context.Context, request *nwv1alpha1.DeleteIPPoolRequest) (*nwv1alpha1.DeleteIPPoolResponse, error) {
	if err := business.DeleteIPPool(ctx, request.Name); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete ippool: %v", err)
	}
	return &nwv1alpha1.DeleteIPPoolResponse{}, nil
}

func (nm *networkManagement) ListMultusConfigs(ctx context.Context, request *nwv1alpha1.ListMultusConfigsRequest) (*nwv1alpha1.ListMultusConfigsResponse, error) {
	multusConfigs, opts, err := business.ListMultusConfigs(ctx, request.Options)
	if errors.IsResourceExpired(err) {
		return nil, status.Errorf(codes.NotFound, "resource expired: %v", err)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list multus configs: %v", err)
	}
	return &nwv1alpha1.ListMultusConfigsResponse{
		Items:   multusConfigs,
		Options: opts,
	}, nil
}

func (nm *networkManagement) ListSubnets(ctx context.Context, request *nwv1alpha1.ListSubnetsRequest) (*nwv1alpha1.ListSubnetsResponse, error) {
	subnets, opts, err := business.ListSubnets(ctx, request.Options)
	if errors.IsResourceExpired(err) {
		return nil, status.Errorf(codes.NotFound, "resource expired: %v", err)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list subnets: %v", err)
	}
	return &nwv1alpha1.ListSubnetsResponse{
		Items:   subnets,
		Options: opts,
	}, nil
}

func (nm *networkManagement) ListIPPools(ctx context.Context, request *nwv1alpha1.ListIPPoolsRequest) (*nwv1alpha1.ListIPPoolsResponse, error) {
	ippools, opts, err := business.ListIPPools(ctx, request.Options)
	if errors.IsResourceExpired(err) {
		return nil, status.Errorf(codes.NotFound, "resource expired: %v", err)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list ippools: %v", err)
	}
	return &nwv1alpha1.ListIPPoolsResponse{
		Items:   ippools,
		Options: opts,
	}, nil
}
