package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/visitor"
)

const (
	BOOLEAN ObjectType = "OBJ_Boolean"
)

func NewBool(v bool) *BoolObject {
	return &BoolObject{Value: v}
}

type BoolObject struct {
	Value bool
}

func (*BoolObject) Type() ObjectType {
	return BOOLEAN
}

func (o *BoolObject) DebugString() string {
	return fmt.Sprintf("Bool(%t)", o.Value)
}

func (o *BoolObject) GoValue() any {
	return o.Value
}

func (o *BoolObject) CoerceTo(t ObjectType) (Object, bool) {
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

func (o *BoolObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*BoolObject)(nil)

// MARK: EqualityEvaluator

func (o *BoolObject) PerformEqualityCheck(r Object) (bool, error) {
	i, err := CoerceTo(r, BOOLEAN)
	if err != nil {
		return false, err
	}

	return o.Value == i.(*BoolObject).Value, nil
}

var _ EqualityEvaluator = (*BoolObject)(nil)
