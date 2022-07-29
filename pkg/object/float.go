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

func (o *Floating) Description() string {
	return fmt.Sprintf("%f", o.Value)
}

func (o Floating) GoValue() any {
	return o.Value
}

func (o *Floating) Cast(t ObjectType) (Object, bool) {
	switch t {
	case INTEGER:
		return &Integer{Value: int64(o.Value)}, true
	case FLOAT:
		return o, true
	case BOOLEAN:
		return CastChain(o, INTEGER, BOOLEAN)
	default:
		return &Void{}, false
	}
}

// MARK: Arithmetic

func (o *Floating) Add(v math.ArithmeticAddition) (math.ArithmeticAddition, bool) {
	return nil, false
}

func (o *Floating) Sub(v math.ArithmeticSubtraction) (math.ArithmeticSubtraction, bool) {
	return nil, false
}

func (o *Floating) Multiply(v math.ArithmeticMultiplication) (math.ArithmeticMultiplication, bool) {
	return nil, false
}

func (o *Floating) Divide(v math.ArithmeticDivision) (math.ArithmeticDivision, bool) {
	return nil, false
}

var _ math.ArithmeticBasic = (*Floating)(nil)

var _ Object = (*Floating)(nil)
