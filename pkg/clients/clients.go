package clients

import (
	"encoding/json"
	"fmt"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubevm.io/vink/pkg/k8s/apis/vink/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/flowcontrol"
	virtv1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/kubecli"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

var Instance = &clients{}

type clients struct {
	kubecli.KubevirtClient

	VinkRestClient    *rest.RESTClient
	KubeOVNRestClient *rest.RESTClient
	KubeRestClient    *rest.RESTClient
}

func InitClients(args ...string) error {
	kubeconfig := GetK8sConfigConfigWithFile(args...)

	kubeconfig.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(100, 200)

	vinkRestClient, err := newRestClientFromRESTConfig(kubeconfig, &v1alpha1.GroupVersion)
	if err != nil {
		return err
	}
	Instance.VinkRestClient = vinkRestClient

	kubeovnRestClient, err := newRestClientFromRESTConfig(kubeconfig, &kubeovnv1.SchemeGroupVersion)
	if err != nil {
		return err
	}
	Instance.KubeOVNRestClient = kubeovnRestClient

	kubevirtClient, err := kubecli.GetKubevirtClientFromRESTConfig(kubeconfig)
	if err != nil {
		return err
	}
	Instance.KubevirtClient = kubevirtClient

	kubeRestClient, err := newRestClientFromRESTConfig(kubeconfig, &corev1.SchemeGroupVersion)
	// kubeRestClient, err := rest.RESTClientFor(kubeconfig)
	if err != nil {
		return err
	}
	Instance.KubeRestClient = kubeRestClient

	return nil
}

func newRestClientFromRESTConfig(kubeconfig *rest.Config, gv *schema.GroupVersion) (*rest.RESTClient, error) {
	shallowCopy := *kubeconfig
	if len(gv.Group) == 0 {
		shallowCopy.APIPath = "/api"
	} else {
		shallowCopy.APIPath = "/apis"
	}
	shallowCopy.GroupVersion = gv
	shallowCopy.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	return rest.RESTClientFor(&shallowCopy)
}

func FromUnstructuredList[T any](obj *unstructured.UnstructuredList) (*T, error) {
	typedObj := new(T)
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), typedObj); err != nil {
		return nil, err
	}

	return typedObj, nil
}

func FromUnstructured[T any](obj *unstructured.Unstructured) (*T, error) {
	typedObj := new(T)
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), typedObj); err != nil {
		return nil, err
	}

	return typedObj, nil
}

func GetGVK(obj runtime.Object) (schema.GroupVersionKind, error) {
	gvks, _, _ := scheme.Scheme.ObjectKinds(obj)
	if len(gvks) < 1 {
		return schema.GroupVersionKind{}, fmt.Errorf("no gvk found")
	}
	return gvks[0], nil
}

func Unstructured[T runtime.Object](obj T) (*unstructured.Unstructured, error) {
	gvk, err := GetGVK(obj)
	if err != nil {
		return nil, fmt.Errorf("no gvk found")
	}
	un := &unstructured.Unstructured{}
	c, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}

	un.SetUnstructuredContent(c)
	un.SetAPIVersion(gvk.GroupVersion().String())
	un.SetKind(gvk.Kind)
	return un, nil
}

func UnstructuredToJSON(obj *unstructured.Unstructured) (string, error) {
	jsonBytes, err := json.Marshal(obj.Object)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
func JSONToUnstructured(crd string) (*unstructured.Unstructured, error) {
	obj := map[string]interface{}{}
	if err := json.Unmarshal([]byte(crd), &obj); err != nil {
		return nil, err
	}

	un := &unstructured.Unstructured{}
	un.SetUnstructuredContent(obj)
	return un, nil
}

func InterfaceToUnstructured(obj any) (*unstructured.Unstructured, error) {
	un := &unstructured.Unstructured{}
	c, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}

	un.SetUnstructuredContent(c)
	return un, nil
}

func InterfaceToObjectMeta(obj any) (*metav1.ObjectMeta, error) {
	un := &unstructured.Unstructured{}
	c, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	un.SetUnstructuredContent(c)

	return &metav1.ObjectMeta{
		Name:                       un.GetName(),
		GenerateName:               un.GetGenerateName(),
		Namespace:                  un.GetNamespace(),
		Labels:                     un.GetLabels(),
		Annotations:                un.GetAnnotations(),
		UID:                        un.GetUID(),
		CreationTimestamp:          un.GetCreationTimestamp(),
		DeletionTimestamp:          un.GetDeletionTimestamp(),
		DeletionGracePeriodSeconds: un.GetDeletionGracePeriodSeconds(),
		Finalizers:                 un.GetFinalizers(),
		OwnerReferences:            un.GetOwnerReferences(),
		ResourceVersion:            un.GetResourceVersion(),
		SelfLink:                   un.GetSelfLink(),
		Generation:                 un.GetGeneration(),
	}, nil
}

func CRDToJSON(obj any) (string, error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func InterfaceToJSON(obj any) (string, error) {
	var un *unstructured.Unstructured
	var err error
	switch payload := obj.(type) {
	case *virtv1.VirtualMachine:
		un, err = Unstructured(payload)
	case *v1alpha1.VirtualMachineSummary:
		un, err = Unstructured(payload)
	case *cdiv1beta1.DataVolume:
		un, err = Unstructured(payload)
	case *corev1.Event:
		un, err = Unstructured(payload)
	default:
		return "", fmt.Errorf("unsupported payload type %T", payload)
	}
	jsonBytes, err := UnstructuredToJSON(un)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), err
}

func JSONToCRD[T runtime.Object](crd string) (T, error) {
	var obj T
	return obj, json.Unmarshal([]byte(crd), &obj)
}
