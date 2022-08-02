package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/visitor"
)

const (
	RETURN ObjectType = "OBJ_Return"
)

func NewReturn(o Object) *ReturnObject {
	return &ReturnObject{Value: o}
}

type ReturnObject struct {
	Value Object
}

func (*ReturnObject) Type() ObjectType {
	return RETURN
}

func (o *ReturnObject) DebugString() string {
	return fmt.Sprintf("Return(%s)", o.Value.DebugString())
}

func (o *ReturnObject) GoValue() any {
	return o.Value.GoValue()
}

func (o *ReturnObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*ReturnObject)(nil)

type ReturnError struct {
	Value   Object
	Message string
}

func (e ReturnError) Error() string {
	return fmt.Sprintf("ReturnError: {%s} %s", e.Value.DebugString(), e.Message)
}
