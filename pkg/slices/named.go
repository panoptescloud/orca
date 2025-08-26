package slices

import (
	"slices"
)

type namedElement interface {
	GetName() string
}

func GetNamedElementIndex[T namedElement](items []T, name string) int {
	for i, item := range items {
		if item.GetName() == name {
			return i
		}
	}

	return -1
}

func NamedElementExists[T namedElement](items []T, name string) bool {
	return GetNamedElementIndex(items, name) != -1
}

func GetNamedElement[T namedElement](items []T, name string) *T {
	index := GetNamedElementIndex(items, name)

	if index == -1 {
		return nil
	}

	return &items[index]
}

func UpsertNamedElement[T namedElement](items []T, item T) []T {
	idx := GetNamedElementIndex(items, item.GetName())

	if idx == -1 {
		return append(items, item)
	}

	left := items[0:idx]
	right := items[(idx + 1):]

	return slices.Concat(left, []T{item}, right)

}
