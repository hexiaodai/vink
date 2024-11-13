package utils

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strings"

// 	"github.com/kubevm.io/vink/apis/types"
// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	"k8s.io/apimachinery/pkg/runtime"
// 	virtv1 "kubevirt.io/api/core/v1"
// )

// // Matches checks if the given object matches the field selector criteria.
// func Matches(payload interface{}, selector *types.FieldSelector) (bool, error) {
// 	var valueMap map[string]interface{}

// 	switch p := payload.(type) {
// 	case *virtv1.VirtualMachine:
// 		valueMap = map[string]interface{}{
// 			"metadata": StructToMap(p.ObjectMeta),
// 			"spec":     StructToMap(p.Spec),
// 			"status":   StructToMap(p.Status),
// 		}
// 	case *virtv1.VirtualMachineInstance:
// 		valueMap = map[string]interface{}{
// 			"metadata": p.ObjectMeta,
// 			"spec":     p.Spec,
// 			"status":   p.Status,
// 		}
// 	default:
// 		return false, fmt.Errorf("unsupported payload type %T", payload)
// 	}

// 	for _, condition := range selector.Conditions {
// 		value, found := nestedFieldValue(valueMap, condition.Fields...)
// 		if !found {
// 			if selector.Operator == types.FieldSelector_OR {
// 				continue
// 			}
// 			return false, nil
// 		}

// 		valueStr := fmt.Sprintf("%v", value)
// 		var match bool
// 		switch condition.Operator {
// 		case types.Condition_EQUAL:
// 			match = valueStr == condition.Value
// 		case types.Condition_NOT_EQUAL:
// 			match = valueStr != condition.Value
// 		case types.Condition_FUZZY:
// 			match = strings.Contains(valueStr, condition.Value)
// 		default:
// 			return false, fmt.Errorf("unsupported operator: %v", condition.Operator)
// 		}

// 		if selector.Operator == types.FieldSelector_AND && !match {
// 			return false, nil
// 		}
// 		if selector.Operator == types.FieldSelector_OR && match {
// 			return true, nil
// 		}
// 	}

// 	return selector.Operator == types.FieldSelector_AND, nil
// }

// // nestedFieldValue 从嵌套的 map 中递归查找字段值
// func nestedFieldValue(valueMap map[string]interface{}, fields ...string) (interface{}, bool) {
// 	current := valueMap
// 	for i, field := range fields {
// 		if val, ok := current[field]; ok {
// 			if i == len(fields)-1 {
// 				return val, true
// 			}
// 			if nested, ok := val.(map[string]interface{}); ok {
// 				current = nested
// 			} else {
// 				return nil, false
// 			}
// 		} else {
// 			return nil, false
// 		}
// 	}
// 	return nil, false
// }
// func StructToMap(v interface{}) map[string]interface{} {
// 	b, _ := json.Marshal(v)
// 	var result map[string]interface{}
// 	json.Unmarshal(b, &result)
// 	return result
// }

// func MatchesUnstructured(payload interface{}, selector *types.FieldSelector) (bool, error) {
// 	payloadMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(payload)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to unmarshal payload to map[string]interface{}: %v", err)
// 	}

// 	for _, condition := range selector.Conditions {
// 		value, found, err := unstructured.NestedFieldNoCopy(payloadMap, condition.Fields...)
// 		if err != nil {
// 			return false, fmt.Errorf("failed to get field %v from payload: %v", condition.Fields, err)
// 		} else if !found {
// 			if selector.Operator == types.FieldSelector_OR {
// 				// If it's an OR operation, not finding the field means the condition doesn't match, but we should continue checking other conditions.
// 				continue
// 			}
// 			return false, nil // If it's an AND operation and field is not found, we return false immediately.
// 		}

// 		valueStr := fmt.Sprintf("%v", value)
// 		var match bool
// 		switch condition.Operator {
// 		case types.Condition_EQUAL:
// 			match = valueStr == condition.Value
// 		case types.Condition_NOT_EQUAL:
// 			match = valueStr != condition.Value
// 		case types.Condition_FUZZY:
// 			match = strings.Contains(valueStr, condition.Value)
// 		default:
// 			return false, fmt.Errorf("unsupported operator: %v", condition.Operator)
// 		}

// 		if selector.Operator == types.FieldSelector_AND && !match {
// 			// For AND operation, if any condition does not match, return false immediately.
// 			return false, nil
// 		}
// 		if selector.Operator == types.FieldSelector_OR && match {
// 			// For OR operation, if any condition matches, return true immediately.
// 			return true, nil
// 		}
// 	}

// 	// If AND operation and all conditions passed, return true.
// 	// If OR operation and no conditions matched, return false.
// 	return selector.Operator == types.FieldSelector_AND, nil
// }
