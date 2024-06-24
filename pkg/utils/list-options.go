package utils

import (
	"github.com/kubevm.io/vink/apis/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConvertToK8sListOptions(in *common.ListOptions) (opts metav1.ListOptions) {
	if in == nil {
		return
	}

	opts = metav1.ListOptions{
		Continue:      in.Continue,
		Limit:         int64(in.Limit),
		LabelSelector: in.LabelsSelector,
		FieldSelector: in.FieldSelector,
	}
	return
}

func ConvertToAPIListOptions(old *common.ListOptions, listMeta metav1.ListMeta) *common.ListOptions {
	output := &common.ListOptions{}
	if old != nil {
		output = old.DeepCopy()
	}
	output.Continue = listMeta.Continue
	return output
}
