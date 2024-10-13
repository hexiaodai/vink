package utils

import (
	"fmt"

	kubeovn "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	apiextensions_v1alpha1 "github.com/kubevm.io/vink/apis/apiextensions/v1alpha1"
	"github.com/kubevm.io/vink/apis/types"
	"github.com/samber/lo"
	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	pkg_types "k8s.io/apimachinery/pkg/types"
	virtv1 "kubevirt.io/api/core/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

// func ConvertGVR(gvri *types.GroupVersionResourceIdentifier) schema.GroupVersionResource {
// 	switch gvri.GetEnum() {
// 	case types.GroupVersionResourceEnum_VIRTUAL_MACHINE:
// 		return gvr.From(virtv1.VirtualMachine{})
// 	case types.GroupVersionResourceEnum_VIRTUAL_MACHINE_INSTANCE:
// 		return gvr.From(virtv1.VirtualMachineInstance{})
// 	case types.GroupVersionResourceEnum_DATA_VOLUME:
// 		return gvr.From(cdiv1.DataVolume{})
// 	case types.GroupVersionResourceEnum_NODE:
// 		return gvr.From(k8sv1.Node{})
// 	case types.GroupVersionResourceEnum_NAMESPACE:
// 		return gvr.From(k8sv1.Namespace{})
// 	case types.GroupVersionResourceEnum_MULTUS:
// 		return gvr.From(netv1.NetworkAttachmentDefinition{})
// 	case types.GroupVersionResourceEnum_SUBNET:
// 		return gvr.From(kubeovn.Subnet{})
// 	case types.GroupVersionResourceEnum_VPC:
// 		return gvr.From(kubeovn.Vpc{})
// 	case types.GroupVersionResourceEnum_IPPOOL:
// 		return gvr.From(kubeovn.IPPool{})
// 	case types.GroupVersionResourceEnum_STORAGE_CLASS:
// 		return gvr.From(storagev1.StorageClass{})
// 	case types.GroupVersionResourceEnum_IPS:
// 		return gvr.From(kubeovn.IP{})
// 	}

// 	if custom := gvri.GetCustom(); custom != nil {
// 		return schema.GroupVersionResource{
// 			Group:    custom.Group,
// 			Version:  custom.Version,
// 			Resource: custom.Resource,
// 		}
// 	}

// 	return schema.GroupVersionResource{}
// }

func ConvertToNamespaceName(payload interface{}) (*types.NamespaceName, error) {
	temp, ok := payload.(*pkg_types.NamespacedName)
	if !ok {
		return nil, fmt.Errorf("unsupported payload type %T", payload)
	}
	return &types.NamespaceName{Namespace: temp.Namespace, Name: temp.Name}, nil
}

func ConvertToCustomResourceDefinition2(payload interface{}) (*apiextensions_v1alpha1.CustomResourceDefinition, error) {
	crd := apiextensions_v1alpha1.CustomResourceDefinition{}

	var metadata metav1.ObjectMeta
	var spec interface{}
	var status interface{}

	switch payload := payload.(type) {
	case *virtv1.VirtualMachine:
		metadata = payload.ObjectMeta
		spec = payload.Spec
		status = payload.Status
	case *virtv1.VirtualMachineInstance:
		metadata = payload.ObjectMeta
		spec = payload.Spec
		status = payload.Status
	case *cdiv1.DataVolume:
		metadata = payload.ObjectMeta
		spec = payload.Spec
		status = payload.Status
	case *k8sv1.Node:
		metadata = payload.ObjectMeta
		spec = payload.Spec
		status = payload.Status
	case *kubeovn.Subnet:
		metadata = payload.ObjectMeta
		spec = payload.Spec
		status = payload.Status
	case *kubeovn.Vpc:
		metadata = payload.ObjectMeta
		spec = payload.Spec
		status = payload.Status
	case *kubeovn.IPPool:
		metadata = payload.ObjectMeta
		spec = payload.Spec
		status = payload.Status
	case *kubeovn.IP:
		metadata = payload.ObjectMeta
		spec = payload.Spec
	default:
		return nil, fmt.Errorf("unsupported payload type %T", payload)
	}

	ownerReferences := []*types.OwnerReference{}
	for _, ownerReference := range metadata.GetOwnerReferences() {
		ownerReferences = append(ownerReferences, &types.OwnerReference{
			ApiVersion:         ownerReference.APIVersion,
			Kind:               ownerReference.Kind,
			Name:               ownerReference.Name,
			Uid:                string(ownerReference.UID),
			Controller:         lo.FromPtr(ownerReference.Controller),
			BlockOwnerDeletion: lo.FromPtr(ownerReference.BlockOwnerDeletion),
		})
	}

	var deletionTimestamp int64
	if metadata.GetDeletionTimestamp() != nil {
		deletionTimestamp = metadata.GetDeletionTimestamp().Time.Unix()
	}

	crd = apiextensions_v1alpha1.CustomResourceDefinition{
		Metadata: &types.ObjectMeta{
			Name:                       metadata.GetName(),
			GenerateName:               metadata.GetGenerateName(),
			Namespace:                  metadata.GetNamespace(),
			Labels:                     metadata.GetLabels(),
			Annotations:                metadata.GetAnnotations(),
			Uid:                        string(metadata.GetUID()),
			CreationTimestamp:          metadata.GetCreationTimestamp().Time.Unix(),
			DeletionTimestamp:          deletionTimestamp,
			DeletionGracePeriodSeconds: lo.FromPtr(metadata.GetDeletionGracePeriodSeconds()),
			Finalizers:                 metadata.GetFinalizers(),
			OwnerReferences:            ownerReferences,
			ResourceVersion:            metadata.GetResourceVersion(),
			SelfLink:                   metadata.GetSelfLink(),
			Generation:                 metadata.GetGeneration(),
		},
	}

	crd.Spec = StructToString(spec)
	crd.Status = StructToString(status)

	return &crd, nil
}

func ConvertToCustomResourceDefinition(payload interface{}) (*apiextensions_v1alpha1.CustomResourceDefinition, error) {
	payloadMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload to map[string]interface{}: %v", err)
	}

	return ConvertUnstructuredToCRD(unstructured.Unstructured{Object: payloadMap})
}

func ConvertUnstructuredToCRD(unStructObj unstructured.Unstructured) (*apiextensions_v1alpha1.CustomResourceDefinition, error) {
	ownerReferences := []*types.OwnerReference{}
	for _, ownerReference := range unStructObj.GetOwnerReferences() {
		ownerReferences = append(ownerReferences, &types.OwnerReference{
			ApiVersion:         ownerReference.APIVersion,
			Kind:               ownerReference.Kind,
			Name:               ownerReference.Name,
			Uid:                string(ownerReference.UID),
			Controller:         lo.FromPtr(ownerReference.Controller),
			BlockOwnerDeletion: lo.FromPtr(ownerReference.BlockOwnerDeletion),
		})
	}

	var deletionTimestamp int64
	if unStructObj.GetDeletionTimestamp() != nil {
		deletionTimestamp = unStructObj.GetDeletionTimestamp().Time.Unix()
	}

	crd := apiextensions_v1alpha1.CustomResourceDefinition{
		Metadata: &types.ObjectMeta{
			Name:                       unStructObj.GetName(),
			GenerateName:               unStructObj.GetGenerateName(),
			Namespace:                  unStructObj.GetNamespace(),
			Labels:                     unStructObj.GetLabels(),
			Annotations:                unStructObj.GetAnnotations(),
			Uid:                        string(unStructObj.GetUID()),
			CreationTimestamp:          unStructObj.GetCreationTimestamp().Time.Unix(),
			DeletionTimestamp:          deletionTimestamp,
			DeletionGracePeriodSeconds: lo.FromPtr(unStructObj.GetDeletionGracePeriodSeconds()),
			Finalizers:                 unStructObj.GetFinalizers(),
			OwnerReferences:            ownerReferences,
			ResourceVersion:            unStructObj.GetResourceVersion(),
			SelfLink:                   unStructObj.GetSelfLink(),
			Generation:                 unStructObj.GetGeneration(),
		},
	}

	spce, _, err := unstructured.NestedFieldNoCopy(unStructObj.Object, "spec")
	// spce, _, err := unstructured.NestedMap(unStructObj.Object, "spec")
	if err != nil {
		return nil, fmt.Errorf("failed to get spce: %v", err)
	}

	status, _, err := unstructured.NestedFieldNoCopy(unStructObj.Object, "status")
	// status, _, err := unstructured.NestedMap(unStructObj.Object, "status")
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %v", err)
	}

	crd.Spec = StructToString(spce)
	crd.Status = StructToString(status)

	return &crd, nil
}

// func getMetadata(obj interface{}) (*types.ObjectMeta, error) {
// 	payloadMap, err := pkg_runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal payload to map[string]interface{}: %v", err)
// 	}

// 	unStructObj := unstructured.Unstructured{Object: payloadMap}

// 	ownerReferences := []*types.OwnerReference{}
// 	for _, ownerReference := range unStructObj.GetOwnerReferences() {
// 		ownerReferences = append(ownerReferences, &types.OwnerReference{
// 			ApiVersion:         ownerReference.APIVersion,
// 			Kind:               ownerReference.Kind,
// 			Name:               ownerReference.Name,
// 			Uid:                string(ownerReference.UID),
// 			Controller:         lo.FromPtr(ownerReference.Controller),
// 			BlockOwnerDeletion: lo.FromPtr(ownerReference.BlockOwnerDeletion),
// 		})
// 	}

// 	var deletionTimestamp int64
// 	if unStructObj.GetDeletionTimestamp() != nil {
// 		deletionTimestamp = unStructObj.GetDeletionTimestamp().Time.Unix()
// 	}

// 	metadata := types.ObjectMeta{
// 		Name:                       unStructObj.GetName(),
// 		GenerateName:               unStructObj.GetGenerateName(),
// 		Namespace:                  unStructObj.GetNamespace(),
// 		Labels:                     unStructObj.GetLabels(),
// 		Annotations:                unStructObj.GetAnnotations(),
// 		Uid:                        string(unStructObj.GetUID()),
// 		CreationTimestamp:          unStructObj.GetCreationTimestamp().Time.Unix(),
// 		DeletionTimestamp:          deletionTimestamp,
// 		DeletionGracePeriodSeconds: lo.FromPtr(unStructObj.GetDeletionGracePeriodSeconds()),
// 		Finalizers:                 unStructObj.GetFinalizers(),
// 		OwnerReferences:            ownerReferences,
// 		ResourceVersion:            unStructObj.GetResourceVersion(),
// 		SelfLink:                   unStructObj.GetSelfLink(),
// 		Generation:                 unStructObj.GetGeneration(),
// 	}

// 	return &metadata, nil
// }

// func ConvertToCustomResourceDefinition(payload interface{}) (*apiextensions_v1alpha1.CustomResourceDefinition, error) {
// 	// return &apiextensions_v1alpha1.CustomResourceDefinition{
// 	// 	Spec: StructToString(payload),
// 	// }, nil

// 	unStructObj := payload.(*virtv1.VirtualMachine)

// 	// payloadMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(payload)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("failed to unmarshal payload to map[string]interface{}: %v", err)
// 	// }

// 	// unStructObj := unstructured.Unstructured{Object: payloadMap}

// 	ownerReferences := []*types.OwnerReference{}
// 	for _, ownerReference := range unStructObj.GetOwnerReferences() {
// 		ownerReferences = append(ownerReferences, &types.OwnerReference{
// 			ApiVersion:         ownerReference.APIVersion,
// 			Kind:               ownerReference.Kind,
// 			Name:               ownerReference.Name,
// 			Uid:                string(ownerReference.UID),
// 			Controller:         lo.FromPtr(ownerReference.Controller),
// 			BlockOwnerDeletion: lo.FromPtr(ownerReference.BlockOwnerDeletion),
// 		})
// 	}

// 	var deletionTimestamp int64
// 	if unStructObj.GetDeletionTimestamp() != nil {
// 		deletionTimestamp = unStructObj.GetDeletionTimestamp().Time.Unix()
// 	}

// 	crd := apiextensions_v1alpha1.CustomResourceDefinition{
// 		Metadata: &types.ObjectMeta{
// 			Name:                       unStructObj.GetName(),
// 			GenerateName:               unStructObj.GetGenerateName(),
// 			Namespace:                  unStructObj.GetNamespace(),
// 			Labels:                     unStructObj.GetLabels(),
// 			Annotations:                unStructObj.GetAnnotations(),
// 			Uid:                        string(unStructObj.GetUID()),
// 			CreationTimestamp:          unStructObj.GetCreationTimestamp().Time.Unix(),
// 			DeletionTimestamp:          deletionTimestamp,
// 			DeletionGracePeriodSeconds: lo.FromPtr(unStructObj.GetDeletionGracePeriodSeconds()),
// 			Finalizers:                 unStructObj.GetFinalizers(),
// 			OwnerReferences:            ownerReferences,
// 			ResourceVersion:            unStructObj.GetResourceVersion(),
// 			SelfLink:                   unStructObj.GetSelfLink(),
// 			Generation:                 unStructObj.GetGeneration(),
// 		},
// 	}

// 	// spce, found, err := unstructured.NestedMap(unStructObj.Object, "spec")
// 	// if !found || err != nil {
// 	// 	return nil, fmt.Errorf("failed to get spce: %v", err)
// 	// }

// 	// status, found, err := unstructured.NestedMap(unStructObj.Object, "status")
// 	// if !found || err != nil {
// 	// 	return nil, fmt.Errorf("failed to get spce: %v", err)
// 	// }

// 	// crd.Spec = StructToString(spce)
// 	// crd.Status = StructToString(status)

// 	crd.Spec = StructToString(unStructObj.Spec)
// 	crd.Status = StructToString(unStructObj.Status)

// 	return &crd, nil
// }
