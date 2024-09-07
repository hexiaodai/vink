package clients

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"fmt"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/flowcontrol"
	"kubevirt.io/client-go/kubecli"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/client-go/dynamic"
)

var _ Clients = (*clients)(nil)

type clients struct {
	dynamicClient  dynamic.Interface
	kubevirtClient kubecli.KubevirtClient
	restClient     *rest.RESTClient
}

func NewClients(args ...string) (Clients, error) {
	cli := clients{}

	kubeconfig := GetK8sConfigConfigWithFile(args...)

	// kubeconfig.APIPath = "/api"
	// kubeconfig.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}
	// kubeconfig.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}

	kubeconfig.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(50, 100)
	dcli, err := dynamic.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	cli.dynamicClient = dcli

	// restClient, err := rest.RESTClientFor(kubeconfig)
	// if err != nil {
	// 	return nil, err
	// }
	// cli.restClient = restClient

	kubevirtClient, err := kubecli.GetKubevirtClientFromRESTConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	cli.kubevirtClient = kubevirtClient

	cli.restClient = kubevirtClient.RestClient()

	return &cli, nil
}

func (cli *clients) GetDynamicKubeClient() dynamic.Interface {
	return cli.dynamicClient
}

func (cli *clients) GetKubeVirtClient() kubecli.KubevirtClient {
	return cli.kubevirtClient
}

func (cli *clients) GetRestClient() *rest.RESTClient {
	return cli.restClient
}

var cli Clients

func GetClients() Clients {
	return cli
}

func InitClients(args ...string) (err error) {
	cli, err = NewClients(args...)
	return
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
