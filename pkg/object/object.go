package object

import "github.com/maddiesch/marble/pkg/object/math"

type Object interface {
	Coercible

	Type() ObjectType

	GoValue() any

	Description() string
}

type BasicArithmeticEvaluator interface {
	PerformBasicArithmeticOperation(math.ArithmeticOperator, Object) (Object, error)
}

type EqualityEvaluator interface {
	PerformEqualityCheck(Object) (bool, error)
}

type ComparisionEvaluator interface {
	EqualityEvaluator

	PerformLessThanComparison(Object) (bool, error)
}

type ConcatingEvaluator interface {
	Concat(Object) (Object, error)
}

type LengthEvaluator interface {
	Len() int
}

type SubscriptEvaluator interface {
	Subscript(Object) (Object, error)
}
