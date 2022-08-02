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

func NewFloat[T constraints.Float](f T) *FloatObject {
	return &FloatObject{Value: float64(f)}
}

type FloatObject struct {
	Value float64
}

func (*FloatObject) Type() ObjectType {
	return FLOAT
}

func (o *FloatObject) DebugString() string {
	return fmt.Sprintf("Float(%f)", o.Value)
}

func (o *FloatObject) GoValue() any {
	return o.Value
}

// TODO: Delete once CastVisitor is complete
func (o *FloatObject) CoerceTo(t ObjectType) (Object, bool) {
	switch t {
	case INTEGER:
		return NewInteger(int64(o.Value)), true
	case FLOAT:
		return o, true
	default:
		return NewVoid(), false
	}
}

func (o *FloatObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*FloatObject)(nil)

// MARK: BasicArithmeticEvaluator

func (o *FloatObject) PerformBasicArithmeticOperation(op math.ArithmeticOperator, val Object) (Object, error) {
	cast, err := CoerceTo(val, FLOAT)
	if err != nil {
		return nil, err
	}
	right := cast.(*FloatObject)

	return NewFloat(math.EvaluateOperation(op, o.Value, right.Value)), nil
}
