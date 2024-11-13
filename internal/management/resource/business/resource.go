package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	resource_v1alpha1 "github.com/kubevm.io/vink/apis/management/resource/v1alpha1"
	"github.com/kubevm.io/vink/apis/types"
	"github.com/kubevm.io/vink/pkg/clients"
	pkg_clients "github.com/kubevm.io/vink/pkg/clients"
	"golang.org/x/sync/errgroup"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	pkg_types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

func List(ctx context.Context, gvr schema.GroupVersionResource, opts *types.ListOptions) ([]string, []*metav1.ObjectMeta, error) {
	cli := clients.Instance.DynamicClient().Resource(gvr)

	items := make([]unstructured.Unstructured, 0)

	switch {
	case len(opts.CustomSelector.NamespaceNames) > 0:
		result, err := listResourcesByCustomNamespaceNames(ctx, cli, opts.CustomSelector.NamespaceNames)
		if err != nil {
			return nil, nil, err
		}
		items = result
	case len(opts.CustomSelector.FieldSelector) > 0:
		result, err := listResourcesByCustomFieldSelector(ctx, cli, opts.Namespace, opts.CustomSelector.FieldSelector)
		if err != nil {
			return nil, nil, err
		}
		items = result
	default:
		result, err := listResourcesByListOptions(ctx, cli, opts)
		if err != nil {
			return nil, nil, err
		}
		items = result
	}

	crds := make([]string, 0, len(items))
	metadatas := make([]*metav1.ObjectMeta, 0, len(items))
	for _, item := range items {
		crd, err := pkg_clients.UnstructuredToJSON(&item)
		// crd, err := utils.ConvertUnstructuredToCRD(item)
		if err != nil {
			return nil, nil, err
		}
		crds = append(crds, crd)
		// FIXME: convert to ObjectMeta
		metadatas = append(metadatas, &metav1.ObjectMeta{
			UID:             item.GetUID(),
			Name:            item.GetName(),
			Namespace:       item.GetNamespace(),
			ResourceVersion: item.GetResourceVersion(),
		})
	}

	return crds, metadatas, nil
}

func listResourcesByListOptions(ctx context.Context, cli dynamic.NamespaceableResourceInterface, opts *types.ListOptions) ([]unstructured.Unstructured, error) {
	res, err := cli.Namespace(opts.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: opts.LabelSelector, FieldSelector: opts.FieldSelector,
		Limit: int64(opts.Limit), Continue: opts.Continue,
	})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}

func listResourcesByCustomFieldSelector(ctx context.Context, cli dynamic.NamespaceableResourceInterface, namespace string, fieldSelector []string) ([]unstructured.Unstructured, error) {
	res, err := cli.Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	filteredItems := make([]unstructured.Unstructured, 0)

	for _, item := range res.Items {
	outerLoop:
		for _, selectorGroup := range fieldSelector {
			conditions := strings.Split(selectorGroup, ",")

			groupMatches := true
			for _, selector := range conditions {
				parts := strings.SplitN(selector, "=", 2)
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid fieldSelector format: %s", selector)
				}

				fieldPath := parts[0]
				expectedValue := parts[1]

				actualValue, found, err := unstructured.NestedString(item.Object, strings.Split(fieldPath, ".")...)
				if err != nil || !found {
					groupMatches = false
					break
				}

				if actualValue != expectedValue {
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

func listResourcesByCustomNamespaceNames(ctx context.Context, cli dynamic.NamespaceableResourceInterface, namespaceNames []*types.NamespaceName) ([]unstructured.Unstructured, error) {
	eg := errgroup.Group{}
	eg.SetLimit(10)

	lock := sync.Mutex{}
	items := make([]unstructured.Unstructured, 0)

	for _, nn := range namespaceNames {
		eg.Go(func() error {
			result, err := cli.Namespace(nn.Namespace).Get(ctx, nn.Name, metav1.GetOptions{})
			if err != nil && !k8serr.IsNotFound(err) {
				return err
			}
			if k8serr.IsNotFound(err) {
				return nil
			}
			lock.Lock()
			items = append(items, *result)
			lock.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return items, nil
}

func Delete(ctx context.Context, gvr schema.GroupVersionResource, nn *types.NamespaceName) error {
	cli := clients.Instance.DynamicClient().Resource(gvr)
	return cli.Namespace(nn.Namespace).Delete(ctx, nn.Name, metav1.DeleteOptions{})
}

func Create(ctx context.Context, gvr schema.GroupVersionResource, crd string) (string, error) {
	// payload := map[string]interface{}{}
	// if err := json.Unmarshal([]byte(lo.FromPtr(crd)), &payload); err != nil {
	// 	return nil, err
	// }

	// obj := unstructured.Unstructured{Object: payload}

	obj, err := pkg_clients.JSONToUnstructured(crd)

	unStructObj, err := clients.Instance.DynamicClient().Resource(gvr).Namespace(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return pkg_clients.UnstructuredToJSON(unStructObj)
	// return utils.ConvertUnstructuredToCRD(*unStructObj)
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

func SendResourceEvent(eventType resource_v1alpha1.EventType, obj interface{}, filterFuncs []FilterFunc, server resource_v1alpha1.ResourceListWatchManagement_ListWatchServer) error {
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

	resp := resource_v1alpha1.ListWatchResponse{
		EventType: eventType,
		Items:     []string{data},
	}
	if err := server.Send(&resp); err != nil {
		return errors.New("failed to send response to client")
	}
	return nil
}
