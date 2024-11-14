package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	resource_v1alpha1 "github.com/kubevm.io/vink/apis/management/resource/v1alpha1"
	"github.com/kubevm.io/vink/apis/types"
	"github.com/kubevm.io/vink/pkg/clients"
	pkg_clients "github.com/kubevm.io/vink/pkg/clients"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	pkg_types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

func List(ctx context.Context, gvr schema.GroupVersionResource, opts *resource_v1alpha1.ListOptions) ([]string, error) {
	cli := clients.Instance.DynamicClient().Resource(gvr)

	items := make([]unstructured.Unstructured, 0)

	switch {
	case len(opts.ArbitraryFieldSelectors) > 0:
		result, err := listResourcesByArbitraryFieldSelectors(ctx, cli, opts.Namespace, opts.ArbitraryFieldSelectors)
		if err != nil {
			return nil, err
		}
		items = result
	default:
		result, err := listResourcesByListOptions(ctx, cli, opts)
		if err != nil {
			return nil, err
		}
		items = result
	}

	crds := make([]string, 0, len(items))
	for _, item := range items {
		crd, err := pkg_clients.UnstructuredToJSON(&item)
		if err != nil {
			return nil, err
		}
		crds = append(crds, crd)
	}

	return crds, nil
}

type arbitraryFieldSelector struct {
	FieldPath     string
	Operator      string
	ExpectedValue string
}

func newArbitraryFieldSelector(selector string) (*arbitraryFieldSelector, error) {
	var operator string
	switch {
	case strings.Contains(selector, "!="):
		operator = "!="
	case strings.Contains(selector, "^="):
		operator = "^="
	case strings.Contains(selector, "$="):
		operator = "$="
	case strings.Contains(selector, "*="):
		operator = "*="
	case strings.Contains(selector, "="):
		operator = "="
	default:
		return nil, fmt.Errorf("unsupported operator in selector: %s", selector)
	}

	parts := strings.SplitN(selector, operator, 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid selector format: %s", selector)
	}

	return &arbitraryFieldSelector{
		FieldPath:     parts[0],
		Operator:      operator,
		ExpectedValue: parts[1],
	}, nil
}

func (fs *arbitraryFieldSelector) matches(item *unstructured.Unstructured) (bool, error) {
	actualValue, found, err := unstructured.NestedString(item.Object, strings.Split(fs.FieldPath, ".")...)
	if err != nil || !found {
		return false, err
	}

	switch fs.Operator {
	case "=":
		return actualValue == fs.ExpectedValue, nil
	case "!=":
		return actualValue != fs.ExpectedValue, nil
	case "^=":
		return strings.HasPrefix(actualValue, fs.ExpectedValue), nil
	case "$=":
		return strings.HasSuffix(actualValue, fs.ExpectedValue), nil
	case "*=":
		return strings.Contains(actualValue, fs.ExpectedValue), nil
	default:
		return false, errors.New("unsupported operator")
	}
}

func listResourcesByListOptions(ctx context.Context, cli dynamic.NamespaceableResourceInterface, opts *resource_v1alpha1.ListOptions) ([]unstructured.Unstructured, error) {
	res, err := cli.Namespace(opts.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: opts.LabelSelector,
		FieldSelector: opts.FieldSelector,
		Limit:         int64(opts.Limit),
		Continue:      opts.Continue,
	})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}

func listResourcesByArbitraryFieldSelectors(ctx context.Context, cli dynamic.NamespaceableResourceInterface, namespace string, fieldSelectors []string) ([]unstructured.Unstructured, error) {
	res, err := cli.Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var filteredItems []unstructured.Unstructured

	for _, item := range res.Items {
	outerLoop:
		for _, selectorGroup := range fieldSelectors {
			conditions := strings.Split(selectorGroup, ",")
			groupMatches := true

			for _, condition := range conditions {
				fieldSelector, err := newArbitraryFieldSelector(condition)
				if err != nil {
					return nil, err
				}

				if match, err := fieldSelector.matches(&item); err != nil || !match {
					groupMatches = false
					break
				}
			}

			if groupMatches {
				filteredItems = append(filteredItems, item)
				break outerLoop
			}
		}
	}

	return filteredItems, nil
}

func Delete(ctx context.Context, gvr schema.GroupVersionResource, nn *types.NamespaceName) error {
	cli := clients.Instance.DynamicClient().Resource(gvr)
	return cli.Namespace(nn.Namespace).Delete(ctx, nn.Name, metav1.DeleteOptions{})
}

func Create(ctx context.Context, gvr schema.GroupVersionResource, crd string) (string, error) {
	obj, err := pkg_clients.JSONToUnstructured(crd)

	unStructObj, err := clients.Instance.DynamicClient().Resource(gvr).Namespace(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return pkg_clients.UnstructuredToJSON(unStructObj)
}

func Update(ctx context.Context, gvr schema.GroupVersionResource, crd string) (string, error) {

	payload := map[string]interface{}{}
	if err := json.Unmarshal([]byte(crd), &payload); err != nil {
		return "", err
	}

	obj := unstructured.Unstructured{Object: payload}

	unStructObj, err := clients.Instance.DynamicClient().Resource(gvr).Namespace(obj.GetNamespace()).Update(ctx, &obj, metav1.UpdateOptions{})
	if err != nil {
		return "", err
	}
	return pkg_clients.UnstructuredToJSON(unStructObj)
}

func Get(ctx context.Context, gvr schema.GroupVersionResource, namespace, name string) (string, error) {
	cli := clients.Instance.DynamicClient().Resource(gvr)
	unStructObj, err := cli.Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return pkg_clients.UnstructuredToJSON(unStructObj)
}

type FilterFunc func(unobj *unstructured.Unstructured) (bool, error)

func trueFilterFunc() FilterFunc {
	return func(_ *unstructured.Unstructured) (bool, error) {
		return true, nil
	}
}

// func FilterFuncWithLabelSelector(opts *resource_v1alpha1.WatchOptions) (FilterFunc, error) {
// 	if len(opts.LabelSelector) == 0 {
// 		return trueFilterFunc(), nil
// 	}
// 	selector, err := labels.Parse(opts.LabelSelector)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse LabelSelector: %v", err)
// 	}

// 	return func(unobj *unstructured.Unstructured) (bool, error) {
// 		return selector.Matches(labels.Set(unobj.GetLabels())), nil
// 	}, nil
// }

func FilterFuncWithFieldSelector(opts *resource_v1alpha1.WatchOptions) (FilterFunc, error) {
	if len(opts.FieldSelector) == 0 {
		return trueFilterFunc(), nil
	}

	return func(unobj *unstructured.Unstructured) (bool, error) {
		for _, selectorGroup := range opts.FieldSelector {
			if len(selectorGroup) == 0 {
				continue
			}
			conditions := strings.Split(selectorGroup, ",")
			groupMatches := true
			for _, condition := range conditions {
				fieldSelector, err := newArbitraryFieldSelector(condition)
				if err != nil {
					return false, fmt.Errorf("failed to parse ArbitraryFieldSelector: %v", err)
				}

				match, err := fieldSelector.matches(unobj)
				if err != nil {
					return false, fmt.Errorf("failed to match field selector: %v", err)
				}
				if !match {
					groupMatches = false
					break
				}
			}
			if groupMatches {
				return true, nil
			}
		}
		return false, nil
	}, nil
}

func DefaultFilterFunc(items []*metav1.ObjectMeta) FilterFunc {
	idx := make(map[string]*metav1.ObjectMeta, len(items))
	for _, metadata := range items {
		ns := pkg_types.NamespacedName{Namespace: metadata.Namespace, Name: metadata.Name}
		idx[ns.String()] = metadata
	}

	return func(unobj *unstructured.Unstructured) (bool, error) {
		ns := pkg_types.NamespacedName{Namespace: unobj.GetNamespace(), Name: unobj.GetName()}
		original, ok := idx[ns.String()]
		if !ok {
			return false, nil
		}

		if unobj.GetUID() != original.GetUID() {
			return true, nil
		}

		originalVersion, err := strconv.Atoi(original.ResourceVersion)
		if err != nil {

			return false, err
		}
		version, err := strconv.Atoi(unobj.GetResourceVersion())
		if err != nil {
			return false, err
		}

		return version > originalVersion, nil
	}
}

func SendResourceEvent(eventType resource_v1alpha1.EventType, obj interface{}, filterFuncs []FilterFunc, server resource_v1alpha1.ResourceWatchManagement_WatchServer) error {
	unobj, err := clients.InterfaceToUnstructured(obj)
	if err != nil {
		return err
	}

	for _, fn := range filterFuncs {
		ok, err := fn(unobj)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}

	data, err := clients.InterfaceToJSON(obj)
	if err != nil {
		return err
	}

	resp := resource_v1alpha1.WatchResponse{
		EventType: eventType,
		Items:     []string{data},
	}
	if err := server.Send(&resp); err != nil {
		return errors.New("failed to send response to client")
	}
	return nil
}
