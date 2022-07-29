package conformance

type Comparable interface {
	Equateable

	LessThan(any) (bool, bool)
}

func CompareLessThan(lhs Comparable, rhs any) (bool, bool) {
	return lhs.LessThan(rhs)
}

func CompareLessThanOrEqual(lhs Comparable, rhs any) (bool, bool) {
	if less, ok := CompareLessThan(lhs, rhs); ok {
		if less {
			return true, true
		}
		return lhs.Equal(rhs)
	}
	return false, false
}

func CompareGreaterThan(lhs Comparable, rhs any) (bool, bool) {
	if less, ok := lhs.LessThan(rhs); ok {
		return !less, true
	}
	return false, false
}

func CompareGreaterThanOrEqual(lhs Comparable, rhs any) (bool, bool) {
	if greater, ok := CompareGreaterThan(lhs, rhs); ok {
		if greater {
			return true, true
		}
		return lhs.Equal(rhs)
	}
	return false, false
}
