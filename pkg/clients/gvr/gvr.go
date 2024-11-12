package gvr

import (
	netv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	kubeovn "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubevm.io/vink/apis/types"
	"github.com/kubevm.io/vink/pkg/k8s/apis/vink/v1alpha1"
	spv2beta1 "github.com/spidernet-io/spiderpool/pkg/k8s/apis/spiderpool.spidernet.io/v2beta1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	virtv1 "kubevirt.io/api/core/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

func init() {
	cdiv1beta1.AddToScheme(scheme.Scheme)
	spv2beta1.AddToScheme(scheme.Scheme)
	virtv1.AddToScheme(scheme.Scheme)
	netv1.AddToScheme(scheme.Scheme)
	kubeovn.AddToScheme(scheme.Scheme)
	storagev1.AddToScheme(scheme.Scheme)
	v1alpha1.AddToScheme(scheme.Scheme)
}

func From[T any](o T) schema.GroupVersionResource {
	switch any(o).(type) {
	case cdiv1beta1.DataVolume, *cdiv1beta1.DataVolume:
		return schema.GroupVersionResource{
			Group:    cdiv1beta1.SchemeGroupVersion.Group,
			Version:  cdiv1beta1.SchemeGroupVersion.Version,
			Resource: "datavolumes",
		}
	case corev1.PersistentVolumeClaim, *corev1.PersistentVolumeClaim:
		return schema.GroupVersionResource{
			Group:    corev1.SchemeGroupVersion.Group,
			Version:  corev1.SchemeGroupVersion.Version,
			Resource: "persistentvolumeclaims",
		}
	case corev1.ConfigMap, *corev1.ConfigMap:
		return schema.GroupVersionResource{
			Group:    corev1.SchemeGroupVersion.Group,
			Version:  corev1.SchemeGroupVersion.Version,
			Resource: "configmaps",
		}
	case storagev1.StorageClass, *storagev1.StorageClass:
		return schema.GroupVersionResource{
			Group:    storagev1.SchemeGroupVersion.Group,
			Version:  storagev1.SchemeGroupVersion.Version,
			Resource: "storageclasses",
		}
	case corev1.Secret, *corev1.Secret:
		return schema.GroupVersionResource{
			Group:    corev1.SchemeGroupVersion.Group,
			Version:  corev1.SchemeGroupVersion.Version,
			Resource: "secrets",
		}
	case corev1.Node, *corev1.Node:
		return schema.GroupVersionResource{
			Group:    corev1.SchemeGroupVersion.Group,
			Version:  corev1.SchemeGroupVersion.Version,
			Resource: "nodes",
		}
	case virtv1.VirtualMachine, *virtv1.VirtualMachine:
		return schema.GroupVersionResource{
			Group:    virtv1.SchemeGroupVersion.Group,
			Version:  virtv1.SchemeGroupVersion.Version,
			Resource: "virtualmachines",
		}
	case virtv1.VirtualMachineInstance, *virtv1.VirtualMachineInstance:
		return schema.GroupVersionResource{
			Group:    virtv1.SchemeGroupVersion.Group,
			Version:  virtv1.SchemeGroupVersion.Version,
			Resource: "virtualmachineinstances",
		}
	case corev1.Namespace, *corev1.Namespace:
		return schema.GroupVersionResource{
			Group:    corev1.SchemeGroupVersion.Group,
			Version:  corev1.SchemeGroupVersion.Version,
			Resource: "namespaces",
		}
	case netv1.NetworkAttachmentDefinition, *netv1.NetworkAttachmentDefinition:
		return schema.GroupVersionResource{
			Group:    netv1.SchemeGroupVersion.Group,
			Version:  netv1.SchemeGroupVersion.Version,
			Resource: "network-attachment-definitions",
		}
	case kubeovn.Subnet, *kubeovn.Subnet:
		return kubeovn.SchemeGroupVersion.WithResource("subnets")
	case kubeovn.Vpc, *kubeovn.Vpc:
		return kubeovn.SchemeGroupVersion.WithResource("vpcs")
	case kubeovn.IPPool, *kubeovn.IPPool:
		return kubeovn.SchemeGroupVersion.WithResource("ippools")
	case kubeovn.IP, *kubeovn.IP:
		return kubeovn.SchemeGroupVersion.WithResource("ips")
	case v1alpha1.VirtualMachineSummary, *v1alpha1.VirtualMachineSummary:
		return v1alpha1.VirtualMachineSummaryGVR
	}

	return schema.GroupVersionResource{}
}

func ResolveGVR(rt types.ResourceType) schema.GroupVersionResource {
	switch rt {
	case types.ResourceType_VIRTUAL_MACHINE:
		return From(virtv1.VirtualMachine{})
	case types.ResourceType_VIRTUAL_MACHINE_INSTANCE:
		return From(virtv1.VirtualMachineInstance{})
	case types.ResourceType_DATA_VOLUME:
		return From(cdiv1beta1.DataVolume{})
	case types.ResourceType_NODE:
		return From(corev1.Node{})
	case types.ResourceType_NAMESPACE:
		return From(corev1.Namespace{})
	case types.ResourceType_MULTUS:
		return From(netv1.NetworkAttachmentDefinition{})
	case types.ResourceType_SUBNET:
		return From(kubeovn.Subnet{})
	case types.ResourceType_VPC:
		return From(kubeovn.Vpc{})
	case types.ResourceType_IPPOOL:
		return From(kubeovn.IPPool{})
	case types.ResourceType_STORAGE_CLASS:
		return From(storagev1.StorageClass{})
	case types.ResourceType_IPS:
		return From(kubeovn.IP{})
	case types.ResourceType_VIRTUAL_MACHINE_SUMMARY:
		return From(v1alpha1.VirtualMachineSummary{})
	}

	return schema.GroupVersionResource{}
}

// func GetGVRFromObjectUsingMapper(obj runtime.Object) (schema.GroupVersionResource, error) {
// 	cachedDiscoveryClient := memory.NewMemCacheClient(clients.GetClients().GetDiscoveryClient())

// 	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedDiscoveryClient)

// 	gvk := obj.GetObjectKind().GroupVersionKind()

// 	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
// 	if err != nil {
// 		return schema.GroupVersionResource{}, err
// 	}

// 	return mapping.Resource, nil
// }
