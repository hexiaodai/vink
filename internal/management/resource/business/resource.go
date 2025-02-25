package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	resource_v1alpha1 "github.com/kubevm.io/vink/apis/management/resource/v1alpha1"
	"github.com/kubevm.io/vink/apis/types"
	"github.com/kubevm.io/vink/pkg/clients"
	pkg_clients "github.com/kubevm.io/vink/pkg/clients"
	"github.com/yalp/jsonpath"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	pkg_types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

func List(ctx context.Context, gvr schema.GroupVersionResource, opts *resource_v1alpha1.ListOptions) ([]string, error) {
	cli := clients.Clients.DynamicClient().Resource(gvr)

	items := make([]unstructured.Unstructured, 0, 0)

	switch {
	case opts.FieldSelectorGroup != nil && len(opts.FieldSelectorGroup.FieldSelectors) > 0:
		result, err := listResourcesByArbitraryFieldSelectors(ctx, cli, opts.Namespace, opts.FieldSelectorGroup)
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
	JsonPath       string
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
		if len(fs.JsonPath) > 0 {
			if value, err := getJsonPathValue(actualValue, fs.JsonPath); err == nil {
				actualValue = value
			}
		}
		actualValue = strings.ToLower(actualValue)

		var expectedValue string
		if len(fs.ExpectedValues) > 0 {
			expectedValue = strings.ToLower(fs.ExpectedValues[0])
		}
		switch fs.Operator {
		// Equals
		case "=":
			if actualValue == expectedValue {
				return true, nil
			}
		// Not equals
		case "!=":
			if actualValue != expectedValue {
				return true, nil
			}
		// Starts with
		case "^=":
			if strings.HasPrefix(actualValue, expectedValue) {
				return true, nil
			}
		// Ends with
		case "$=":
			if strings.HasSuffix(actualValue, expectedValue) {
				return true, nil
			}
		// Contains
		case "*=":
			if strings.Contains(actualValue, expectedValue) {
				return true, nil
			}
		// Matches one of
		case "~=":
			for _, ev := range fs.ExpectedValues {
				ev = strings.ToLower(ev)
				if actualValue == strings.TrimSpace(ev) {
					return true, nil
				}
			}
		// Does not match any of
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

var errInvalidJsonString = errors.New("invalid JSON string")

func getJsonPathValue(input string, jsonPath string) (string, error) {
	var jsonObject interface{}
	if err := json.Unmarshal([]byte(input), &jsonObject); err != nil {
		return "", errInvalidJsonString
	}

	value, err := jsonpath.Read(jsonObject, jsonPath)
	if err != nil {
		return "", err
	}

	var valueStr string
	switch v := value.(type) {
	case string:
		valueStr = v
	case float64, int, bool:
		valueStr = fmt.Sprintf("%v", v)
	default:
		return "", fmt.Errorf("unsupported JSONPath value type: %T", v)
	}

	return valueStr, nil
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

func listResourcesByArbitraryFieldSelectors(ctx context.Context, cli dynamic.NamespaceableResourceInterface, namespace string, fieldSelectorGroup *types.FieldSelectorGroup) ([]unstructured.Unstructured, error) {
	res, err := cli.Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	filterFunc, err := FilterFuncWithFieldSelector(fieldSelectorGroup)
	if err != nil {
		return nil, err
	}

	var filteredItems []unstructured.Unstructured

	for _, item := range res.Items {
		if ok, err := filterFunc(&item); err != nil {
			return nil, err
		} else if ok {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems, nil
}

func Delete(ctx context.Context, gvr schema.GroupVersionResource, nn *types.NamespaceName) error {
	cli := clients.Clients.DynamicClient().Resource(gvr)
	return cli.Namespace(nn.Namespace).Delete(ctx, nn.Name, metav1.DeleteOptions{})
}

func Create(ctx context.Context, gvr schema.GroupVersionResource, crd string) (string, error) {
	obj, err := pkg_clients.JSONToUnstructured(crd)
	if err != nil {
		return "", err
	}

	unStructObj, err := clients.Clients.DynamicClient().Resource(gvr).Namespace(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
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

	unStructObj, err := clients.Clients.DynamicClient().Resource(gvr).Namespace(obj.GetNamespace()).Update(ctx, &obj, metav1.UpdateOptions{})
	if err != nil {
		return "", err
	}
	return pkg_clients.UnstructuredToJSON(unStructObj)
}

func Get(ctx context.Context, gvr schema.GroupVersionResource, namespace, name string) (string, error) {
	cli := clients.Clients.DynamicClient().Resource(gvr)
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

func FilterFuncWithFieldSelector(fieldSelectorGroup *types.FieldSelectorGroup) (FilterFunc, error) {
	if fieldSelectorGroup == nil || len(fieldSelectorGroup.FieldSelectors) == 0 {
		return trueFilterFunc(), nil
	}

	return func(unobj *unstructured.Unstructured) (bool, error) {
		for _, selector := range fieldSelectorGroup.FieldSelectors {
			fieldSelector := &arbitraryFieldSelector{
				FieldPath:      selector.FieldPath,
				JsonPath:       selector.JsonPath,
				Operator:       selector.Operator,
				ExpectedValues: selector.Values,
			}

			match, err := fieldSelector.matches(unobj)
			if err != nil {
				return false, fmt.Errorf("failed to match field selector: %v", err)
			}
			switch fieldSelectorGroup.Operator {
			case "||":
				if match {
					return true, nil
				}
			case "&&", "":
				if !match {
					return false, nil
				}
			default:
				return false, fmt.Errorf("unknown field selector group operator: %s", fieldSelectorGroup.Operator)
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

func SendReadyEventOnce(readyOnce *sync.Once, server resource_v1alpha1.ResourceWatchManagement_WatchServer) error {
	var err error
	readyOnce.Do(func() {
		err = server.Send(&resource_v1alpha1.WatchResponse{EventType: resource_v1alpha1.EventType_READY})
	})
	if err != nil {
		return fmt.Errorf("failed to send readiness status to client: %w", err)
	}
	return nil
}

func SendResourceEventWithReady(readyOnce *sync.Once, eventType resource_v1alpha1.EventType, obj interface{}, filterFuncs []FilterFunc, server resource_v1alpha1.ResourceWatchManagement_WatchServer) error {
	if err := SendReadyEventOnce(readyOnce, server); err != nil {
		return err
	}

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
		return fmt.Errorf("failed to send response to client: %w", err)
	}

	return nil
}
