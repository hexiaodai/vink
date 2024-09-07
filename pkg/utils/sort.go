package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func ExtractFieldValue(obj interface{}, fields []string) (interface{}, error) {
	payloadMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload to map[string]interface{}: %v", err)
	}

	value, _, err := unstructured.NestedFieldNoCopy(payloadMap, fields...)
	if err != nil {
		return nil, fmt.Errorf("failed to get fields %v: %v", fields, err)
	}

	if reflect.DeepEqual(fields, []string{"metadata", "creationTimestamp"}) {
		var timestamp metav1.Time
		timestamp.UnmarshalQueryParameter(value.(string))
		return timestamp.Time, nil
	}

	return value, nil
}

func ExtractNumber(s string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindString(s)
	if matches == "" {
		return 0, fmt.Errorf("no number found")
	}
	return strconv.Atoi(matches)
}

func SortFunc(items []interface{}, fields []string, order string) func(i, j int) bool {
	return func(i, j int) bool {
		valA, errA := ExtractFieldValue(items[i], fields)
		valB, errB := ExtractFieldValue(items[j], fields)

		if errA != nil || errB != nil {
			return false
		}

		switch vA := valA.(type) {
		case string:
			vB, ok := valB.(string)
			if !ok {
				return false
			}
			numA, errA := ExtractNumber(vA)
			numB, errB := ExtractNumber(vB)
			if errA == nil && errB == nil {
				if order == "asc" {
					return numA < numB
				}
				return numA > numB
			}
			if order == "asc" {
				return vA < vB
			}
			return vA > vB
		case int, float64:
			vB, ok := valB.(float64)
			if !ok {
				return false
			}
			if order == "asc" {
				return vA.(float64) < vB
			}
			return vA.(float64) > vB
		case time.Time:
			vB, ok := valB.(time.Time)
			if !ok {
				return false
			}
			if order == "asc" {
				return vA.Before(vB)
			}
			return vA.After(vB)
		default:
			return false
		}
	}
}
