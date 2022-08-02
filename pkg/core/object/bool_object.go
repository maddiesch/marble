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

func (o *BoolObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*BoolObject)(nil)

// MARK: EqualityEvaluator

func (o *BoolObject) PerformEqualityCheck(r Object) (bool, error) {
	boolean, err := CastObjectTo(r, BOOLEAN)
	if err != nil {
		panic(err) // TODO: Better Error handling
	}

	return o.Value == boolean.(*BoolObject).Value, nil
}

var _ EqualityEvaluator = (*BoolObject)(nil)
