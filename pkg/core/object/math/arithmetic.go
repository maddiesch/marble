package math

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

type ArithmeticOperator rune

const (
	OperationAdd      ArithmeticOperator = '+'
	OperationSubtract ArithmeticOperator = '-'
	OperationMultiply ArithmeticOperator = '*'
	OperationDivide   ArithmeticOperator = '/'
)

func OperatorFor(s string) (ArithmeticOperator, bool) {
	var v ArithmeticOperator
	r := []rune(s)
	if len(r) == 0 {
		return v, false
	}

	switch r[0] {
	case '+':
		return OperationAdd, true
	case '-':
		return OperationSubtract, true
	case '*':
		return OperationMultiply, true
	case '/':
		return OperationDivide, true
	}

	return v, false
}

func EvaluateOperation[T constraints.Integer | constraints.Float](op ArithmeticOperator, lhs, rhs T) T {
	switch op {
	case OperationAdd:
		return lhs + rhs
	case OperationSubtract:
		return lhs - rhs
	case OperationMultiply:
		return lhs * rhs
	case OperationDivide:
		return lhs / rhs
	default:
		panic("unable to perform arithmetic operation for operator: " + strconv.QuoteRune(rune(op)))
	}
}
