package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/object/math"
	"github.com/maddiesch/marble/pkg/core/visitor"
	"golang.org/x/exp/constraints"
)

const (
	INTEGER ObjectType = "OBJ_Integer"
)

func NewInteger[T constraints.Integer](i T) *IntegerObject {
	return &IntegerObject{Value: int64(i)}
}

type IntegerObject struct {
	Value int64
}

func (*IntegerObject) Type() ObjectType {
	return INTEGER
}

func (o *IntegerObject) DebugString() string {
	return fmt.Sprintf("Int(%d)", o.Value)
}

func (o *IntegerObject) GoValue() any {
	return o.Value
}

func (o *IntegerObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*IntegerObject)(nil)

// MARK: BasicArithmeticEvaluator

func (o *IntegerObject) PerformBasicArithmeticOperation(op math.ArithmeticOperator, val Object) (Object, error) {
	integer, err := CastObjectTo(val, INTEGER)
	if err != nil {
		panic(err) // TODO: Better error handling
	}

	return NewInteger(math.EvaluateOperation(op, o.Value, integer.(*IntegerObject).Value)), nil
}

var _ BasicArithmeticEvaluator = (*IntegerObject)(nil)

// MARK: ComparisionEvaluator

func (o *IntegerObject) PerformEqualityCheck(r Object) (bool, error) {
	integer, err := CastObjectTo(r, INTEGER)
	if err != nil {
		panic(err) // TODO: Better error handling
	}

	return o.Value == integer.(*IntegerObject).Value, nil
}

func (o *IntegerObject) PerformLessThanComparison(r Object) (bool, error) {
	integer, err := CastObjectTo(r, INTEGER)
	if err != nil {
		panic(err) // TODO: Better error handling
	}

	return o.Value < integer.(*IntegerObject).Value, nil
}

var _ ComparisionEvaluator = (*IntegerObject)(nil)
