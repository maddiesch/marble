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

func (o *IntegerObject) CoerceTo(t ObjectType) (Object, bool) {
	switch t {
	case INTEGER:
		return o, true
	case FLOAT:
		return NewFloat(float64(o.Value)), true
	case BOOLEAN:
		return Bool(o.Value > 0), true
	default:
		return NewVoid(), false
	}
}

func (o *IntegerObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*IntegerObject)(nil)

// MARK: BasicArithmeticEvaluator

func (o *IntegerObject) PerformBasicArithmeticOperation(op math.ArithmeticOperator, val Object) (Object, error) {
	cast, err := CoerceTo(val, INTEGER)
	if err != nil {
		return nil, err
	}
	right := cast.(*IntegerObject)

	return NewInteger(math.EvaluateOperation(op, o.Value, right.Value)), nil
}

var _ BasicArithmeticEvaluator = (*IntegerObject)(nil)

// MARK: ComparisionEvaluator

func (o *IntegerObject) PerformEqualityCheck(r Object) (bool, error) {
	i, err := CoerceTo(r, INTEGER)
	if err != nil {
		return false, err
	}

	return o.Value == i.(*IntegerObject).Value, nil
}

func (o *IntegerObject) PerformLessThanComparison(r Object) (bool, error) {
	i, err := CoerceTo(r, INTEGER)
	if err != nil {
		return false, err
	}

	return o.Value < i.(*IntegerObject).Value, nil
}

var _ ComparisionEvaluator = (*IntegerObject)(nil)
