package conformance

type Equateable interface {
	Equal(any) (bool, bool)
}

func Equal(lhs Equateable, rhs any) (bool, bool) {
	return lhs.Equal(rhs)
}

func NotEqual(lhs Equateable, rhs any) (bool, bool) {
	if eq, ok := Equal(lhs, rhs); ok {
		return !eq, true
	}
	return false, false
}
