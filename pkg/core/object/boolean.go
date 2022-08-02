package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/visitor"
)

const (
	BOOLEAN ObjectType = "OBJ_Boolean"
)

func Bool(v bool) *Boolean {
	return &Boolean{Value: v}
}

type Boolean struct {
	Value bool
}

func (o *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (o *Boolean) DebugString() string {
	return fmt.Sprintf("Bool(%t)", o.Value)
}

func (o *Boolean) GoValue() any {
	return o.Value
}

func (o *Boolean) CoerceTo(t ObjectType) (Object, bool) {
	switch t {
	case BOOLEAN:
		return o, true
	case INTEGER:
		if o.Value == true {
			return NewInteger(1), true
		} else {
			return NewInteger(0), true
		}
	case FLOAT:
		if o.Value == true {
			return NewFloat(1.0), true
		} else {
			return NewFloat(0.0), true
		}
	default:
		return NewVoid(), false
	}
}

func (o *Boolean) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*Boolean)(nil)

// MARK: EqualityEvaluator

func (o *Boolean) PerformEqualityCheck(r Object) (bool, error) {
	i, err := CoerceTo(r, BOOLEAN)
	if err != nil {
		return false, err
	}

	return o.Value == i.(*Boolean).Value, nil
}

var _ EqualityEvaluator = (*Boolean)(nil)
