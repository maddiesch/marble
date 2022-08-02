package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/visitor"
)

const (
	STRING ObjectType = "OBJ_String"
)

func NewString(s string) *StringObject {
	return &StringObject{Value: s}
}

type StringObject struct {
	Value string
}

func (*StringObject) Type() ObjectType {
	return STRING
}

func (o *StringObject) DebugString() string {
	return fmt.Sprintf("String(%s)", o.Value)
}

func (o *StringObject) GoValue() any {
	return o.Value
}

func (o *StringObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*StringObject)(nil)
