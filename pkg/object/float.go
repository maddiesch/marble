package object

import (
	"fmt"

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

func (o *Floating) CoerceTo(t ObjectType) (Object, bool) {
	switch t {
	case INTEGER:
		return Int(int64(o.Value)), true
	case FLOAT:
		return o, true
	case BOOLEAN:
		o, err := ChainCoerceTo(o, INTEGER, BOOLEAN)
		if err != nil {
			return nil, false
		}
		return o, true
	default:
		return &Void{}, false
	}
}

var _ Object = (*Floating)(nil)
