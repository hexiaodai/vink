package clients

import (
	"encoding/json"
	"fmt"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubevm.io/vink/pkg/k8s/apis/vink/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/flowcontrol"
	"kubevirt.io/client-go/kubecli"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	virtv1 "kubevirt.io/api/core/v1"
)

var _ Clients = (*clients)(nil)

type clients struct {
	dynamicClient     dynamic.Interface
	kubevirtClient    kubecli.KubevirtClient
	discoveryClient   discovery.DiscoveryInterface
	k8sConfig         *rest.Config
	vinkRestClient    *rest.RESTClient
	kubeovnRestClient *rest.RESTClient
}

func NewClients(args ...string) (Clients, error) {
	cli := clients{}

	kubeconfig := GetK8sConfigConfigWithFile(args...)

	kubeconfig.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(100, 200)

	cli.k8sConfig = kubeconfig

	dcli, err := dynamic.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	cli.dynamicClient = dcli

	vinkRestClient, err := vinkRestClientFromRESTConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	cli.vinkRestClient = vinkRestClient

	kubeovnRestClient, err := kubeovnRestClientFromRESTConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	cli.kubeovnRestClient = kubeovnRestClient

	kubevirtClient, err := kubecli.GetKubevirtClientFromRESTConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	cli.kubevirtClient = kubevirtClient

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create discovery client: %v", err)
	}
	cli.discoveryClient = discoveryClient

	return &cli, nil
}

func vinkRestClientFromRESTConfig(kubeconfig *rest.Config) (*rest.RESTClient, error) {
	shallowCopy := *kubeconfig
	shallowCopy.APIPath = "/apis"
	shallowCopy.GroupVersion = &v1alpha1.GroupVersion
	shallowCopy.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	return rest.RESTClientFor(&shallowCopy)
}

func kubeovnRestClientFromRESTConfig(kubeconfig *rest.Config) (*rest.RESTClient, error) {
	shallowCopy := *kubeconfig
	shallowCopy.APIPath = "/apis"
	shallowCopy.GroupVersion = &kubeovnv1.SchemeGroupVersion
	shallowCopy.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	return rest.RESTClientFor(&shallowCopy)
}

func (cli *clients) GetVinkRestClient() *rest.RESTClient {
	return cli.vinkRestClient
}

func (cli *clients) GetKubeovnRestClient() *rest.RESTClient {
	return cli.kubeovnRestClient
}

func (cli *clients) GetDynamicKubeClient() dynamic.Interface {
	return cli.dynamicClient
}

func (cli *clients) GetKubeVirtClient() kubecli.KubevirtClient {
	return cli.kubevirtClient
}

func (cli *clients) GetDiscoveryClient() discovery.DiscoveryInterface {
	return cli.discoveryClient
}

func (cli *clients) GetKubeConfig() *rest.Config {
	return cli.k8sConfig
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
