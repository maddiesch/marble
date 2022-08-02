package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/object/math"
	"golang.org/x/exp/constraints"
)

const (
	FLOAT ObjectType = "OBJ_Float"
)

func Float[T constraints.Float](f T) *Floating {
	return &Floating{Value: float64(f)}
}

type Floating struct {
	Value float64
}

func (*Floating) Type() ObjectType {
	return FLOAT
}

func (o *Floating) DebugString() string {
	return fmt.Sprintf("Float(%f)", o.Value)
}

func (o Floating) GoValue() any {
	return o.Value
}

func (o *Floating) CoerceTo(t ObjectType) (Object, bool) {
	switch t {
	case INTEGER:
		return Int(int64(o.Value)), true
	case FLOAT:
		return o, true
	default:
		return &Void{}, false
	}
}

var _ Object = (*Floating)(nil)

// MARK: BasicArithmeticEvaluator

func (o *Floating) PerformBasicArithmeticOperation(op math.ArithmeticOperator, val Object) (Object, error) {
	cast, err := CoerceTo(val, FLOAT)
	if err != nil {
		return nil, err
	}
	right := cast.(*Floating)

	return Float(math.EvaluateOperation(op, o.Value, right.Value)), nil
}
