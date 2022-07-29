package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/object/conformance"
	"github.com/maddiesch/marble/pkg/object/math"
	"golang.org/x/exp/constraints"
)

const (
	INTEGER ObjectType = "OBJ_Integer"
)

func Int[T constraints.Integer](i T) *Integer {
	return &Integer{Value: int64(i)}
}

type Integer struct {
	Value int64
}

func (*Integer) Type() ObjectType {
	return INTEGER
}

func (o *Integer) Description() string {
	return fmt.Sprintf("Int(%d)", o.Value)
}

func (o *Integer) GoValue() any {
	return o.Value
}

func (o *Integer) CoerceTo(t ObjectType) (Object, bool) {
	switch t {
	case INTEGER:
		return o, true
	case FLOAT:
		return Float(float64(o.Value)), true
	case BOOLEAN:
		return Bool(o.Value > 0), true
	default:
		return &Void{}, false
	}
}

var _ Object = (*Integer)(nil)

// MARK: Comparable

func (o *Integer) Equal(v any) (bool, bool) {
	switch v := v.(type) {
	case *Integer:
		return o.Value == v.Value, true
	case Coercible:
		if cast, ok := o.CoerceTo(o.Type()); ok {
			return o.Equal(cast)
		} else {
			return false, false
		}
	default:
		return false, false
	}
}

func (o *Integer) LessThan(v any) (bool, bool) {
	switch v := v.(type) {
	case *Integer:
		return o.Value < v.Value, true
	case Coercible:
		if cast, ok := o.CoerceTo(o.Type()); ok {
			return o.Equal(cast)
		} else {
			return false, false
		}
	default:
		return false, false
	}
}

var _ conformance.Comparable = (*Integer)(nil)

// MARK: BasicArithmeticEvaluator

func (o *Integer) PerformBasicArithmeticOperation(op math.ArithmeticOperator, val Object) (Object, error) {
	cast, err := CoerceTo(val, INTEGER)
	if err != nil {
		return nil, err
	}
	right := cast.(*Integer)

	return Int(math.EvaluateOperation(op, o.Value, right.Value)), nil
}

var _ BasicArithmeticEvaluator = (*Integer)(nil)
