package util

import (
	"fmt"
	"strings"
)

func JoinList[T any](list []T, sep string) string {
	elem := make([]string, len(list))
	for i, it := range list {
		elem[i] = fmt.Sprint(it)
	}
	return strings.Join(elem, sep)
}

func RemoveFromList[T any](list []T, index int) []T {
	newList := []T{}
	for i, ele := range list {
		if i == index {
			// Remove this element
			continue
		}
		newList = append(newList, ele)
	}

	return newList
}
