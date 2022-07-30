package stack

type Stack[T any] struct {
	list []T
}

func New[T any](v ...T) *Stack[T] {
	return &Stack[T]{list: v}
}

func (s *Stack[T]) Len() int { return len(s.list) }

func (s *Stack[T]) UFirst() T {
	return s.list[0]
}

func (s *Stack[T]) First() (T, bool) {
	if len(s.list) == 0 {
		var e T
		return e, false
	}

	return s.UFirst(), true
}

func (s *Stack[T]) ULast() T {
	return s.list[len(s.list)-1]
}

func (s *Stack[T]) Last() (T, bool) {
	if len(s.list) == 0 {
		var e T
		return e, false
	}

	return s.ULast(), true
}

func (s *Stack[T]) Push(v T) {
	s.list = append(s.list, v)
}

func (s *Stack[T]) Pop() T {
	if len(s.list) == 0 {
		panic("attempting to pop an empty stack")
	}

	val := s.list[len(s.list)-1]
	s.list = s.list[:len(s.list)-1]
	return val
}

func (s *Stack[T]) RevIter() []T {
	i := s.Iter()
	reverse(i)
	return i
}

func (s *Stack[T]) Iter() []T {
	return s.copy()
}

func (s *Stack[T]) At(i int) T {
	return s.list[i]
}

func (s *Stack[T]) Search(fn func(T) bool) (int, T) {
	for i, v := range s.list {
		if fn(v) {
			return i, v
		}
	}
	var v T
	return -1, v
}

func (s *Stack[T]) copy() []T {
	i := make([]T, len(s.list))
	copy(i, s.list)
	return i
}

func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
