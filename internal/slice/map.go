package slice

func Map[R any, T any](s []T, fn func(T) R) []R {
	return MapI(s, func(_ int, v T) R {
		return fn(v)
	})
}

func MapI[R any, T any](s []T, fn func(int, T) R) []R {
	r := make([]R, len(s))
	for i, v := range s {
		r[i] = fn(i, v)
	}
	return r
}
