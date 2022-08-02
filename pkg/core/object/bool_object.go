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
