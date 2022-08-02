package object

import (
	"fmt"
	"strings"

	"github.com/maddiesch/marble/internal/collection"
	"github.com/maddiesch/marble/pkg/core/visitor"
)

const (
	ARRAY ObjectType = "OBJ_Array"
)

func NewArray(v []Object) *ArrayObject {
	return &ArrayObject{Elements: v}
}

type ArrayObject struct {
	Elements []Object
}

func (*ArrayObject) Type() ObjectType {
	return ARRAY
}

func (o *ArrayObject) DebugString() string {
	elements := collection.MapSlice(o.Elements, func(e Object) string {
		return e.DebugString()
	})
	return fmt.Sprintf("Array(%s)", strings.Join(elements, ", "))
}

func (o *ArrayObject) GoValue() any {
	return o.Elements
}

func (o *ArrayObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*ArrayObject)(nil)
