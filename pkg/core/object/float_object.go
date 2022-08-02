package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/visitor"
	"golang.org/x/exp/constraints"
)

const (
	FLOAT ObjectType = "OBJ_Float"
)

func NewFloat[T constraints.Float](f T) *FloatObject {
	return &FloatObject{Value: float64(f)}
}

type FloatObject struct {
	Value float64
}

func (*FloatObject) Type() ObjectType {
	return FLOAT
}

func (o *FloatObject) DebugString() string {
	return fmt.Sprintf("Float(%f)", o.Value)
}

func (o *FloatObject) GoValue() any {
	return o.Value
}

func (o *FloatObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*FloatObject)(nil)
