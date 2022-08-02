package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/visitor"
	"golang.org/x/exp/constraints"
)

const (
	INTEGER ObjectType = "OBJ_Integer"
)

func NewInteger[T constraints.Integer](i T) *IntegerObject {
	return &IntegerObject{Value: int64(i)}
}

type IntegerObject struct {
	Value int64
}

func (*IntegerObject) Type() ObjectType {
	return INTEGER
}

func (o *IntegerObject) DebugString() string {
	return fmt.Sprintf("Int(%d)", o.Value)
}

func (o *IntegerObject) GoValue() any {
	return o.Value
}

func (o *IntegerObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*IntegerObject)(nil)
