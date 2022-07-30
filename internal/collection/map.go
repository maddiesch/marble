package collection

func MapSlice[R any, T any](s []T, fn func(T) R) []R {
	return MapSliceI(s, func(_ int, v T) R {
		return fn(v)
	})
}

func MapSliceI[R any, T any](s []T, fn func(int, T) R) []R {
	r := make([]R, len(s))
	for i, v := range s {
		r[i] = fn(i, v)
	}
	return r
}

func MapMap[R any, K comparable, V any](m map[K]V, fn func(K, V) R) []R {
	r := make([]R, 0, len(m))
	for k, v := range m {
		r = append(r, fn(k, v))
	}
	return r
}
