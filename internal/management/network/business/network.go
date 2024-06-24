package business

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubevm.io/vink/apis/common"
	nwv1alpha1 "github.com/kubevm.io/vink/apis/management/network/v1alpha1"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/clients/gvr"
	"github.com/kubevm.io/vink/pkg/utils"
	spv2beta1 "github.com/spidernet-io/spiderpool/pkg/k8s/apis/spiderpool.spidernet.io/v2beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func ListNodesNetworkInterfaces(ctx context.Context, opts *common.ListOptions) ([]*nwv1alpha1.NodeNetworkInterface, *common.ListOptions, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()
	unobj, err := dcli.Resource(gvr.From(corev1.Node{})).List(ctx, utils.ConvertToK8sListOptions(opts))
	if err != nil {
		return nil, nil, err
	}

	obj, err := clients.FromUnstructuredList[corev1.NodeList](unobj)
	if err != nil {
		return nil, nil, err
	}
	nodes := make([]*nwv1alpha1.NodeNetworkInterface, 0, len(obj.Items))
	for _, item := range obj.Items {
		node, err := crdToAPINodeNetworkInterface(&item)
		if err != nil {
			return nil, nil, err
		}
		nodes = append(nodes, node)
	}

	return nodes, utils.ConvertToAPIListOptions(opts, obj.ListMeta), nil
}

func ListMultusConfigs(ctx context.Context, opts *common.ListOptions) ([]*nwv1alpha1.MultusConfig, *common.ListOptions, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()
	unobj, err := dcli.Resource(gvr.From(spv2beta1.SpiderMultusConfig{})).Namespace(metav1.NamespaceSystem).List(ctx, utils.ConvertToK8sListOptions(opts))
	if err != nil {
		return nil, nil, err
	}

	obj, err := clients.FromUnstructuredList[spv2beta1.SpiderMultusConfigList](unobj)
	if err != nil {
		return nil, nil, err
	}

	multusConfigs := make([]*nwv1alpha1.MultusConfig, 0, len(obj.Items))
	for _, item := range obj.Items {
		mcfg, err := crdToAPIMultusConfig(&item)
		if err != nil {
			return nil, nil, err
		}
		multusConfigs = append(multusConfigs, mcfg)
	}

	return multusConfigs, utils.ConvertToAPIListOptions(opts, obj.ListMeta), nil
}

func ListSubnets(ctx context.Context, opts *common.ListOptions) ([]*nwv1alpha1.Subnet, *common.ListOptions, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()
	unobj, err := dcli.Resource(gvr.From(spv2beta1.SpiderSubnet{})).List(ctx, utils.ConvertToK8sListOptions(opts))
	if err != nil {
		return nil, nil, err
	}

	obj, err := clients.FromUnstructuredList[spv2beta1.SpiderSubnetList](unobj)
	if err != nil {
		return nil, nil, err
	}

	subnets := make([]*nwv1alpha1.Subnet, 0, len(obj.Items))
	for _, item := range obj.Items {
		subnet, err := crdToAPISubnet(&item)
		if err != nil {
			return nil, nil, err
		}
		subnets = append(subnets, subnet)
	}
	return subnets, utils.ConvertToAPIListOptions(opts, obj.ListMeta), nil
}

func ListIPPools(ctx context.Context, opts *common.ListOptions) ([]*nwv1alpha1.IPPool, *common.ListOptions, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()
	unobj, err := dcli.Resource(gvr.From(spv2beta1.SpiderIPPool{})).List(ctx, utils.ConvertToK8sListOptions(opts))
	if err != nil {
		return nil, nil, err
	}

	obj, err := clients.FromUnstructuredList[spv2beta1.SpiderIPPoolList](unobj)
	if err != nil {
		return nil, nil, err
	}

	ippools := make([]*nwv1alpha1.IPPool, 0, len(obj.Items))
	for _, item := range obj.Items {
		ippool, err := crdToAPIIPPool(&item)
		if err != nil {
			return nil, nil, err
		}
		ippools = append(ippools, ippool)
	}
	return ippools, utils.ConvertToAPIListOptions(opts, obj.ListMeta), nil
}

func CreateMultusConfig(ctx context.Context, name, nic string) (*nwv1alpha1.MultusConfig, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()

	cfg := generateMacvlanCNIConfigCRD(name, nic)
	un, _ := clients.Unstructured(cfg)
	unObj, err := dcli.Resource(gvr.From(spv2beta1.SpiderMultusConfig{})).Namespace(metav1.NamespaceSystem).Create(ctx, un, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	obj, err := clients.FromUnstructured[spv2beta1.SpiderMultusConfig](unObj)
	if err != nil {
		return nil, err
	}
	return crdToAPIMultusConfig(obj)
}

func CreateSubnet(ctx context.Context, name string, config *nwv1alpha1.SubnetConfig) (*nwv1alpha1.Subnet, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()

	subnet := generateSubnetCRD(name, config)
	un, _ := clients.Unstructured(subnet)
	unObj, err := dcli.Resource(gvr.From(spv2beta1.SpiderSubnet{})).Create(ctx, un, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	obj, err := clients.FromUnstructured[spv2beta1.SpiderSubnet](unObj)
	if err != nil {
		return nil, err
	}
	return crdToAPISubnet(obj)
}

func CreateIPPool(ctx context.Context, name string, config *nwv1alpha1.IPPoolConfig) (*nwv1alpha1.IPPool, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()

	ippool := generateIPPoolCRD(name, config)
	un, _ := clients.Unstructured(ippool)
	unObj, err := dcli.Resource(gvr.From(spv2beta1.SpiderIPPool{})).Create(ctx, un, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	obj, err := clients.FromUnstructured[spv2beta1.SpiderIPPool](unObj)
	if err != nil {
		return nil, err
	}
	return crdToAPIIPPool(obj)
}

func DeleteMultusConfig(ctx context.Context, name string) error {
	dcli := clients.GetClients().GetDynamicKubeClient()
	err := dcli.Resource(gvr.From(spv2beta1.SpiderMultusConfig{})).Namespace(metav1.NamespaceSystem).Delete(ctx, name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func DeleteSubnet(ctx context.Context, name string) error {
	dcli := clients.GetClients().GetDynamicKubeClient()
	err := dcli.Resource(gvr.From(spv2beta1.SpiderSubnet{})).Delete(ctx, name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func DeleteIPPool(ctx context.Context, name string) error {
	dcli := clients.GetClients().GetDynamicKubeClient()
	err := dcli.Resource(gvr.From(spv2beta1.SpiderIPPool{})).Delete(ctx, name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}
