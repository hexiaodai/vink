package gvr

import (
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
	case spv2beta1.SpiderMultusConfig, *spv2beta1.SpiderMultusConfig:
		return schema.GroupVersionResource{
			Group:    spv2beta1.SchemeGroupVersion.Group,
			Version:  spv2beta1.SchemeGroupVersion.Version,
			Resource: "spidermultusconfigs",
		}
	case spv2beta1.SpiderSubnet, *spv2beta1.SpiderSubnet:
		return schema.GroupVersionResource{
			Group:    spv2beta1.SchemeGroupVersion.Group,
			Version:  spv2beta1.SchemeGroupVersion.Version,
			Resource: "spidersubnets",
		}
	case spv2beta1.SpiderIPPool, *spv2beta1.SpiderIPPool:
		return schema.GroupVersionResource{
			Group:    spv2beta1.SchemeGroupVersion.Group,
			Version:  spv2beta1.SchemeGroupVersion.Version,
			Resource: "spiderippools",
		}
	case spv2beta1.SpiderEndpoint, *spv2beta1.SpiderEndpoint:
		return schema.GroupVersionResource{
			Group:    spv2beta1.SchemeGroupVersion.Group,
			Version:  spv2beta1.SchemeGroupVersion.Version,
			Resource: "spiderendpoints",
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
	}
	return schema.GroupVersionResource{}
}
