package business

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nwv1alpha1 "github.com/kubevm.io/vink/apis/management/network/v1alpha1"
	"github.com/kubevm.io/vink/pkg/utils"
	"github.com/samber/lo"
	spv2beta1 "github.com/spidernet-io/spiderpool/pkg/k8s/apis/spiderpool.spidernet.io/v2beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
	corev1 "k8s.io/api/core/v1"
)

func crdToAPINodeNetworkInterface(in *corev1.Node) (*nwv1alpha1.NodeNetworkInterface, error) {
	inter := nwv1alpha1.NodeNetworkInterface{
		Node: &nwv1alpha1.NodeNetworkInterface_Node{
			Name: in.Name,
		},
	}
	// if value, ok := in.Annotations[annotation.IoVinkNodeNetworkInterface.Name]; ok {
	// 	items := []*node.NetworkInterface{}
	// 	if err := json.Unmarshal([]byte(value), &items); err != nil {
	// 		return nil, err
	// 	}
	// 	for _, item := range items {
	// 		inter.NetworkInterface = append(inter.NetworkInterface, &nwv1alpha1.NodeNetworkInterface_NetworkInterface{
	// 			Name:    item.Name,
	// 			State:   item.State,
	// 			Ip:      item.IP,
	// 			Subnet:  item.Subnet,
	// 			Gateway: item.Gateway,
	// 		})
	// 	}
	// }
	return &inter, nil
}

func crdToAPIMultusConfig(in *spv2beta1.SpiderMultusConfig) (*nwv1alpha1.MultusConfig, error) {
	pbSpec, err := utils.ConvertToProtoStruct(in.Spec)
	if err != nil {
		return nil, err
	}
	return &nwv1alpha1.MultusConfig{
		Name:              in.Name,
		Spec:              pbSpec,
		CreationTimestamp: timestamppb.New(in.CreationTimestamp.Time),
	}, nil
}

func crdToAPISubnet(in *spv2beta1.SpiderSubnet) (*nwv1alpha1.Subnet, error) {
	pbSpec, err := utils.ConvertToProtoStruct(in.Spec)
	if err != nil {
		return nil, err
	}
	pbStatus, err := utils.ConvertToProtoStruct(in.Status)
	if err != nil {
		return nil, err
	}
	return &nwv1alpha1.Subnet{
		Name:              in.Name,
		Spec:              pbSpec,
		Status:            pbStatus,
		CreationTimestamp: timestamppb.New(in.CreationTimestamp.Time),
	}, nil
}

func crdToAPIIPPool(in *spv2beta1.SpiderIPPool) (*nwv1alpha1.IPPool, error) {
	pbSpec, err := utils.ConvertToProtoStruct(in.Spec)
	if err != nil {
		return nil, err
	}
	pbStatus, err := utils.ConvertToProtoStruct(in.Status)
	if err != nil {
		return nil, err
	}
	return &nwv1alpha1.IPPool{
		Name:              in.Name,
		Spec:              pbSpec,
		Status:            pbStatus,
		CreationTimestamp: timestamppb.New(in.CreationTimestamp.Time),
	}, nil
}

func generateMacvlanCNIConfigCRD(name, nic string) *spv2beta1.SpiderMultusConfig {
	return &spv2beta1.SpiderMultusConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: metav1.NamespaceSystem,
		},
		Spec: spv2beta1.MultusCNIConfigSpec{
			CniType: lo.ToPtr("macvlan"),
			MacvlanConfig: &spv2beta1.SpiderMacvlanCniConfig{
				Master: []string{nic},
			},
		},
	}
}

func generateSubnetCRD(name string, cfg *nwv1alpha1.SubnetConfig) *spv2beta1.SpiderSubnet {
	subnet := spv2beta1.SpiderSubnet{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: spv2beta1.SubnetSpec{
			Gateway:    &cfg.Gateway,
			IPVersion:  lo.ToPtr[int64](4),
			IPs:        cfg.Ips,
			Subnet:     cfg.Subnet,
			ExcludeIPs: cfg.ExcludeIps,
		},
	}
	rs := []spv2beta1.Route{}
	for _, r := range cfg.Routes {
		rs = append(rs, spv2beta1.Route{Dst: r.Dst, Gw: r.Gw})
	}
	subnet.Spec.Routes = rs

	return &subnet
}

func generateIPPoolCRD(name string, cfg *nwv1alpha1.IPPoolConfig) *spv2beta1.SpiderIPPool {
	ippool := spv2beta1.SpiderIPPool{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: spv2beta1.IPPoolSpec{
			Gateway:    &cfg.Gateway,
			IPVersion:  lo.ToPtr[int64](4),
			IPs:        cfg.Ips,
			Subnet:     cfg.Subnet,
			ExcludeIPs: cfg.ExcludeIps,
		},
	}
	rs := []spv2beta1.Route{}
	for _, r := range cfg.Routes {
		rs = append(rs, spv2beta1.Route{Dst: r.Dst, Gw: r.Gw})
	}
	ippool.Spec.Routes = rs

	return &ippool
}
