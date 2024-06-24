package business

import (
	"context"

	"github.com/kubevm.io/vink/apis/common"
	nsv1alpha1 "github.com/kubevm.io/vink/apis/management/namespace/v1alpha1"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/clients/gvr"
	"github.com/kubevm.io/vink/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListNamespaces(ctx context.Context, opts *common.ListOptions) ([]*nsv1alpha1.Namespace, *common.ListOptions, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()
	unobj, err := dcli.Resource(gvr.From(corev1.Namespace{})).List(ctx, utils.ConvertToK8sListOptions(opts))
	if err != nil {
		return nil, nil, err
	}
	obj, err := clients.FromUnstructuredList[corev1.NamespaceList](unobj)
	if err != nil {
		return nil, nil, err
	}
	namespaces := make([]*nsv1alpha1.Namespace, 0, len(obj.Items))
	for _, item := range obj.Items {
		dv, err := crdToAPINamespace(&item)
		if err != nil {
			return nil, nil, err
		}
		namespaces = append(namespaces, dv)
	}

	return namespaces, utils.ConvertToAPIListOptions(opts, obj.ListMeta), nil
}

func CreateNamespace(ctx context.Context, name string, config *nsv1alpha1.NamespaceConfig) (*nsv1alpha1.Namespace, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()

	dvcrd, err := generateNamespaceCRD(name, config)
	if err != nil {
		return nil, err
	}
	un, _ := clients.Unstructured(dvcrd)
	unObj, err := dcli.Resource(gvr.From(corev1.Namespace{})).Create(ctx, un, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	obj, err := clients.FromUnstructured[corev1.Namespace](unObj)
	if err != nil {
		return nil, err
	}

	return crdToAPINamespace(obj)
}

func DeleteNamespace(ctx context.Context, name string) error {
	dcli := clients.GetClients().GetDynamicKubeClient()

	err := dcli.Resource(gvr.From(corev1.Namespace{})).Delete(ctx, name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}

	return err
}
