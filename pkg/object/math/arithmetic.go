package math

type ArithmeticAddition interface {
	Add(ArithmeticAddition) (ArithmeticAddition, bool)
}

type ArithmeticSubtraction interface {
	Sub(ArithmeticSubtraction) (ArithmeticSubtraction, bool)
}

type ArithmeticMultiplication interface {
	Multiply(ArithmeticMultiplication) (ArithmeticMultiplication, bool)
}

type ArithmeticDivision interface {
	Divide(ArithmeticDivision) (ArithmeticDivision, bool)
}

type ArithmeticModulus interface {
	Mod(ArithmeticModulus) (ArithmeticModulus, bool)
}

type ArithmeticPower interface {
	Power(ArithmeticPower) (ArithmeticPower, bool)
}

type ArithmeticBasic interface {
	ArithmeticAddition
	ArithmeticSubtraction
	ArithmeticMultiplication
	ArithmeticDivision
}

func Add[T ArithmeticAddition](lhs T, rhs T) (T, bool) {
	if r, ok := lhs.Add(rhs); ok {
		if tr, ok := r.(T); ok {
			return tr, true
		}
	}
	var z T
	return z, false
}

func Subtract[T ArithmeticSubtraction](lhs T, rhs T) (T, bool) {
	if r, ok := lhs.Sub(rhs); ok {
		if tr, ok := r.(T); ok {
			return tr, true
		}
	}
	var z T
	return z, false
}

func Multiply[T ArithmeticMultiplication](lhs T, rhs T) (T, bool) {
	if r, ok := lhs.Multiply(rhs); ok {
		if tr, ok := r.(T); ok {
			return tr, true
		}
	}
	var z T
	return z, false
}

func Divide[T ArithmeticDivision](lhs T, rhs T) (T, bool) {
	if r, ok := lhs.Divide(rhs); ok {
		if tr, ok := r.(T); ok {
			return tr, true
		}
	}
	var z T
	return z, false
}

func Mod[T ArithmeticModulus](lhs T, rhs T) (T, bool) {
	if r, ok := lhs.Mod(rhs); ok {
		if tr, ok := r.(T); ok {
			return tr, true
		}
	}
	var z T
	return z, false
}

func Pow[T ArithmeticPower](lhs T, rhs T) (T, bool) {
	if r, ok := lhs.Power(rhs); ok {
		if tr, ok := r.(T); ok {
			return tr, true
		}
	}
	var z T
	return z, false
}
