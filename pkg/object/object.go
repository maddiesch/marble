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
