package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
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
	FieldPath      string
	Operator       string
	ExpectedValues []string
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
	case strings.Contains(selector, "~="):
		operator = "~="
	case strings.Contains(selector, "!~="):
		operator = "!~="
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
		FieldPath: parts[0],
		Operator:  operator,
		// ExpectedValue: parts[1],
	}, nil
}

func parseFieldPath(fieldPath string) []string {
	re := regexp.MustCompile(`\\\.`)
	parsedField := re.ReplaceAllString(fieldPath, "\u0000")
	parts := strings.Split(parsedField, ".")

	for i := range parts {
		parts[i] = strings.ReplaceAll(parts[i], "\u0000", ".")
	}
	return parts
}

func (fs *arbitraryFieldSelector) matches(item *unstructured.Unstructured) (bool, error) {
	fieldPathParts := parseFieldPath(fs.FieldPath)
	actualValues, err := getValuesFromFieldPath(item.Object, fieldPathParts)
	if err != nil {
		return false, err
	}

	for _, actualValue := range actualValues {
		actualValue = strings.ToLower(actualValue)
		var expectedValue string
		if len(fs.ExpectedValues) > 0 {
			expectedValue = strings.ToLower(fs.ExpectedValues[0])
		}
		switch fs.Operator {
		case "=":
			if actualValue == expectedValue {
				return true, nil
			}
		case "!=":
			if actualValue != expectedValue {
				return true, nil
			}
		case "^=":
			if strings.HasPrefix(actualValue, expectedValue) {
				return true, nil
			}
		case "$=":
			if strings.HasSuffix(actualValue, expectedValue) {
				return true, nil
			}
		case "*=":
			if strings.Contains(actualValue, expectedValue) {
				return true, nil
			}
		case "~=":
			for _, ev := range fs.ExpectedValues {
				ev = strings.ToLower(ev)
				if actualValue == strings.TrimSpace(ev) {
					return true, nil
				}
			}
		case "!~=":
			for _, ev := range fs.ExpectedValues {
				ev = strings.ToLower(ev)
				if actualValue == strings.TrimSpace(ev) {
					return false, nil
				}
			}
			return true, nil
		default:
			return false, errors.New("unsupported operator")
		}
	}

	// If no value matched, return false
	return false, nil
}

// Helper function to recursively resolve the field path
func getValuesFromFieldPath(obj map[string]interface{}, fieldPathParts []string) ([]string, error) {
	if len(fieldPathParts) == 0 {
		return nil, nil
	}

	currentPart := fieldPathParts[0]
	remainingParts := fieldPathParts[1:]

	// Handle array access
	if strings.HasSuffix(currentPart, "]") {
		openBracketIndex := strings.Index(currentPart, "[")
		if openBracketIndex == -1 {
			return nil, fmt.Errorf("invalid field path: %s", currentPart)
		}

		fieldName := currentPart[:openBracketIndex]
		arrayIndex := currentPart[openBracketIndex+1 : len(currentPart)-1]

		// Access the array
		rawArray, found, err := unstructured.NestedFieldNoCopy(obj, fieldName)
		if err != nil || !found {
			return nil, nil // Treat as not found
		}

		array, ok := rawArray.([]interface{})
		if !ok {
			return nil, fmt.Errorf("field %s is not an array", fieldName)
		}

		// Handle specific index [n]
		if arrayIndex != "*" {
			index, err := strconv.Atoi(arrayIndex)
			if err != nil || index < 0 || index >= len(array) {
				return nil, nil // Treat as not found
			}
			return getValuesFromFieldPath(array[index].(map[string]interface{}), remainingParts)
		}

		// Handle wildcard [*]
		var results []string
		for _, item := range array {
			switch value := item.(type) {
			case map[string]interface{}:
				values, err := getValuesFromFieldPath(value, remainingParts)
				if err != nil {
					continue
				}
				results = append(results, values...)
			case string:
				if len(remainingParts) == 0 {
					results = append(results, value)
				}
			default:
				continue
			}
		}
		return results, nil
	}

	// Handle regular field
	rawValue, found, err := unstructured.NestedFieldNoCopy(obj, currentPart)
	if err != nil || !found {
		return nil, nil // Treat as not found
	}

	// If there are no more parts, return the value as a string
	if len(remainingParts) == 0 {
		if strValue, ok := rawValue.(string); ok {
			return []string{strValue}, nil
		}
		return nil, fmt.Errorf("field %s is not a string", currentPart)
	}

	// Continue to the next part
	subObj, ok := rawValue.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("field %s is not an object", currentPart)
	}
	return getValuesFromFieldPath(subObj, remainingParts)
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

func FilterFuncWithFieldSelector(opts *resource_v1alpha1.WatchOptions) (FilterFunc, error) {
	if opts.FieldSelectorGroup == nil || len(opts.FieldSelectorGroup.FieldSelectors) == 0 {
		return trueFilterFunc(), nil
	}

	return func(unobj *unstructured.Unstructured) (bool, error) {
		for _, selector := range opts.FieldSelectorGroup.FieldSelectors {
			fieldSelector := &arbitraryFieldSelector{
				FieldPath:      selector.FieldPath,
				Operator:       selector.Operator,
				ExpectedValues: selector.Values,
			}

			match, err := fieldSelector.matches(unobj)
			if err != nil {
				return false, fmt.Errorf("failed to match field selector: %v", err)
			}
			switch opts.FieldSelectorGroup.Operator {
			case "||":
				if match {
					return true, nil
				}
			case "&&":
				if !match {
					return false, nil
				}
			default:
				return false, fmt.Errorf("unknown field selector group operator: %s", opts.FieldSelectorGroup.Operator)
			}
		}
		return true, nil
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
