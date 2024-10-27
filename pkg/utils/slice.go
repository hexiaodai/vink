package utils

import (
	"golang.org/x/exp/constraints"
)

func CompareArrays[T constraints.Ordered](arr1, arr2 []T) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	counts := make(map[T]int)

	for _, elem := range arr1 {
		counts[elem]++
	}

	for _, elem := range arr2 {
		if counts[elem] == 0 {
			return false
		}
		counts[elem]--
	}

	for _, count := range counts {
		if count != 0 {
			return false
		}
	}

	return true
}
