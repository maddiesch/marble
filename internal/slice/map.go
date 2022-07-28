package slice

func Map[R any, T any](s []T, fn func(T) R) []R {
	r := make([]R, len(s))
	for i, v := range s {
		r[i] = fn(v)
	}
	return r
}
