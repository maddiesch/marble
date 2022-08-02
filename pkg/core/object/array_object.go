package object

import (
	"fmt"
	"strings"

	"github.com/maddiesch/marble/internal/collection"
	"github.com/maddiesch/marble/pkg/core/evaluator/runtime"
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

// MARK: ComparisionEvaluator

func (s *ArrayObject) Concat(r Object) (Object, error) {
	array, err := CastObjectTo(r, ARRAY)
	if err != nil {
		panic(err) // TODO: Return a valid go error
	}

	return NewArray(append(s.Elements, array.(*ArrayObject).Elements...)), nil
}

var _ ConcatingEvaluator = (*ArrayObject)(nil)

func (o *ArrayObject) Len() int {
	return len(o.Elements)
}

var _ LengthEvaluator = (*ArrayObject)(nil)

func (o *ArrayObject) Subscript(k Object) (Object, error) {
	integer, err := CastObjectTo(k, INTEGER)
	if err != nil {
		return nil, runtime.NewError(runtime.ArgumentError,
			"Unable to access array using given argument, must be Int",
			runtime.ErrorValue("CoercionError", err),
			runtime.ErrorValue("Argument", k),
		)
	}
	return o.Elements[int(integer.(*IntegerObject).Value)], nil
}

var _ SubscriptEvaluator = (*ArrayObject)(nil)
