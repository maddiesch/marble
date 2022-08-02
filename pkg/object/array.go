package object

import (
	"fmt"
	"strings"

	"github.com/maddiesch/marble/internal/collection"
	"github.com/maddiesch/marble/pkg/evaluator/runtime"
)

const (
	ARRAY ObjectType = "OBJ_Array"
)

func Array(v []Object) *ArrayObject {
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

func (o *ArrayObject) CoerceTo(t ObjectType) (Object, bool) {
	switch t {
	case ARRAY:
		return o, true
	default:
		return &Void{}, false
	}
}

var _ Object = (*ArrayObject)(nil)

// MARK: ComparisionEvaluator

func (s *ArrayObject) Concat(r Object) (Object, error) {
	r, err := CoerceTo(r, ARRAY)
	if err != nil {
		return nil, err
	}

	return Array(append(s.Elements, r.(*ArrayObject).Elements...)), nil
}

var _ ConcatingEvaluator = (*ArrayObject)(nil)

func (o *ArrayObject) Len() int {
	return len(o.Elements)
}

var _ LengthEvaluator = (*ArrayObject)(nil)

func (o *ArrayObject) Subscript(k Object) (Object, error) {
	index, err := CoerceTo(k, INTEGER)
	if err != nil {
		return nil, runtime.NewError(runtime.ArgumentError,
			"Unable to access array using given argument, must be Int",
			runtime.ErrorValue("CoercionError", err),
			runtime.ErrorValue("Argument", k),
		)
	}
	return o.Elements[int(index.(*Integer).Value)], nil
}

var _ SubscriptEvaluator = (*ArrayObject)(nil)
