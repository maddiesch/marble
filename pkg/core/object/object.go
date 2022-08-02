package object

import (
	"github.com/maddiesch/marble/pkg/core/object/math"
	"github.com/maddiesch/marble/pkg/debug"
)

type Object interface {
	Coercible
	debug.Description

	Type() ObjectType

	GoValue() any
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
