package business

import (
	"context"

	"github.com/kubevm.io/vink/apis/common"
	dvv1alpha1 "github.com/kubevm.io/vink/apis/management/datavolume/v1alpha1"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/clients/gvr"
	"github.com/kubevm.io/vink/pkg/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

func ListDataVolumes(ctx context.Context, namespace string, opts *common.ListOptions) ([]*dvv1alpha1.DataVolume, *common.ListOptions, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()
	unobj, err := dcli.Resource(gvr.From(cdiv1beta1.DataVolume{})).Namespace(namespace).List(ctx, utils.ConvertToK8sListOptions(opts))
	if err != nil {
		return nil, nil, err
	}
	obj, err := clients.FromUnstructuredList[cdiv1beta1.DataVolumeList](unobj)
	if err != nil {
		return nil, nil, err
	}
	datavolumes := make([]*dvv1alpha1.DataVolume, 0, len(obj.Items))
	for _, item := range obj.Items {
		dv, err := crdToAPIDataVolume(&item)
		if err != nil {
			return nil, nil, err
		}
		datavolumes = append(datavolumes, dv)
	}

	return datavolumes, utils.ConvertToAPIListOptions(opts, obj.ListMeta), nil
}

func CreateDataVolumes(ctx context.Context, namespace, name string, config *dvv1alpha1.DataVolumeConfig) (*dvv1alpha1.DataVolume, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()

	dvcrd, err := generateDataVolumeCRD(namespace, name, config)
	if err != nil {
		return nil, err
	}
	un, _ := clients.Unstructured(dvcrd)
	unObj, err := dcli.Resource(gvr.From(cdiv1beta1.DataVolume{})).Namespace(namespace).Create(ctx, un, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	obj, err := clients.FromUnstructured[cdiv1beta1.DataVolume](unObj)
	if err != nil {
		return nil, err
	}

	return crdToAPIDataVolume(obj)
}

func DeleteDataVolume(ctx context.Context, namespace, name string) error {
	dcli := clients.GetClients().GetDynamicKubeClient()

	err := dcli.Resource(gvr.From(cdiv1beta1.DataVolume{})).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}

	return err
}
