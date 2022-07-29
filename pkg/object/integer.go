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
	return fmt.Sprintf("%d", o.Value)
}

func (o *Integer) GoValue() any {
	return o.Value
}

func (o *Integer) Cast(t ObjectType) (Object, bool) {
	switch t {
	case INTEGER:
		return o, true
	case FLOAT:
		return &Floating{Value: float64(o.Value)}, true
	case BOOLEAN:
		return &Boolean{Value: o.Value > 0}, true
	default:
		return &Void{}, false
	}
}

var _ Object = (*Integer)(nil)

// MARK: Arithmetic

func (o *Integer) Add(v math.ArithmeticAddition) (math.ArithmeticAddition, bool) {
	if c, ok := v.(Castable); ok {
		if i, ok := c.Cast(o.Type()); ok {
			return Int(o.Value + i.(*Integer).Value), true
		}
	}
	return nil, false
}

func (o *Integer) Sub(v math.ArithmeticSubtraction) (math.ArithmeticSubtraction, bool) {
	if c, ok := v.(Castable); ok {
		if i, ok := c.Cast(o.Type()); ok {
			return Int(o.Value - i.(*Integer).Value), true
		}
	}
	return nil, false
}

func (o *Integer) Multiply(v math.ArithmeticMultiplication) (math.ArithmeticMultiplication, bool) {
	if c, ok := v.(Castable); ok {
		if i, ok := c.Cast(o.Type()); ok {
			return Int(o.Value * i.(*Integer).Value), true
		}
	}
	return nil, false
}

func (o *Integer) Divide(v math.ArithmeticDivision) (math.ArithmeticDivision, bool) {
	if c, ok := v.(Castable); ok {
		if i, ok := c.Cast(o.Type()); ok {
			return Int(o.Value / i.(*Integer).Value), true
		}
	}
	return nil, false
}

var _ math.ArithmeticBasic = (*Integer)(nil)

// MARK: Comparable

func (o *Integer) Equal(v any) (bool, bool) {
	switch v := v.(type) {
	case *Integer:
		return o.Value == v.Value, true
	case Castable:
		if cast, ok := o.Cast(o.Type()); ok {
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
	case Castable:
		if cast, ok := o.Cast(o.Type()); ok {
			return o.Equal(cast)
		} else {
			return false, false
		}
	default:
		return false, false
	}
}

var _ conformance.Comparable = (*Integer)(nil)
