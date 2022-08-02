package collection

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func MapSlice[R any, T any](s []T, fn func(T) R) []R {
	return MapSliceI(s, func(_ int, v T) R {
		return fn(v)
	})
}

func MapSliceI[R any, T any](s []T, fn func(int, T) R) []R {
	r, _ := MapSliceErr(s, func(i int, v T) (R, error) {
		return fn(i, v), nil
	})
	return r
}

func MapSliceErr[R any, T any](s []T, fn func(int, T) (R, error)) ([]R, error) {
	r := make([]R, len(s))
	var err error
	for i, v := range s {
		r[i], err = fn(i, v)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func MapMap[R any, K comparable, V any](m map[K]V, fn func(K, V) R) []R {
	r := make([]R, 0, len(m))
	for k, v := range m {
		r = append(r, fn(k, v))
	}
	return r
}

func Keys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))

	for k := range m {
		r = append(r, k)
	}

	return r
}

func SortedKeys[K constraints.Ordered, V any](m map[K]V) []K {
	keys := Keys(m)

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}
