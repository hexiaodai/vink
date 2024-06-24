package business

import (
	nsv1alpha1 "github.com/kubevm.io/vink/apis/management/namespace/v1alpha1"
	"github.com/kubevm.io/vink/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	corev1 "k8s.io/api/core/v1"
)

func crdToAPINamespace(in *corev1.Namespace) (*nsv1alpha1.Namespace, error) {
	namespace := nsv1alpha1.Namespace{
		Name:              in.Name,
		Namespace:         utils.MustConvertToProtoStruct(in),
		CreationTimestamp: timestamppb.New(in.CreationTimestamp.Time),
	}
	return &namespace, nil
}

func generateNamespaceCRD(name string, config *nsv1alpha1.NamespaceConfig) (*corev1.Namespace, error) {
	return nil, nil
}
