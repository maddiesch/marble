package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/object/math"
	"github.com/maddiesch/marble/pkg/core/visitor"
	"golang.org/x/exp/constraints"
)

const (
	FLOAT ObjectType = "OBJ_Float"
)

func NewFloat[T constraints.Float](f T) *FloatingObject {
	return &FloatingObject{Value: float64(f)}
}

type FloatingObject struct {
	Value float64
}

func (*FloatingObject) Type() ObjectType {
	return FLOAT
}

func (o *FloatingObject) DebugString() string {
	return fmt.Sprintf("Float(%f)", o.Value)
}

func (o *FloatingObject) GoValue() any {
	return o.Value
}

func (o *FloatingObject) CoerceTo(t ObjectType) (Object, bool) {
	switch t {
	case INTEGER:
		return NewInteger(int64(o.Value)), true
	case FLOAT:
		return o, true
	default:
		return NewVoid(), false
	}
}

func (o *FloatingObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*FloatingObject)(nil)

// MARK: BasicArithmeticEvaluator

func (o *FloatingObject) PerformBasicArithmeticOperation(op math.ArithmeticOperator, val Object) (Object, error) {
	cast, err := CoerceTo(val, FLOAT)
	if err != nil {
		return nil, err
	}
	right := cast.(*FloatingObject)

	return NewFloat(math.EvaluateOperation(op, o.Value, right.Value)), nil
}
