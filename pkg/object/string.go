package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/object/conformance"
)

const (
	STRING ObjectType = "OBJ_String"
)

func String(s string) *StringObject {
	return &StringObject{Value: s}
}

type StringObject struct {
	Value string
}

func (*StringObject) Type() ObjectType {
	return STRING
}

func (o *StringObject) Description() string {
	return fmt.Sprintf("String(%s)", o.Value)
}

func (o *StringObject) GoValue() any {
	return o.Value
}

func (o *StringObject) CoerceTo(t ObjectType) (Object, bool) {
	switch t {
	case STRING:
		return o, true
	default:
		return &Void{}, false
	}
}

var _ Object = (*StringObject)(nil)

// MARK: Comparable

func (o *StringObject) Equal(v any) (bool, bool) {
	switch v := v.(type) {
	case *StringObject:
		return o.Value == v.Value, true
	case Coercible:
		if cast, ok := o.CoerceTo(o.Type()); ok {
			return o.Equal(cast)
		} else {
			return false, false
		}
	default:
		return false, false
	}
}

func (o *StringObject) LessThan(v any) (bool, bool) {
	switch v := v.(type) {
	case *StringObject:
		return o.Value < v.Value, true
	case Coercible:
		if cast, ok := o.CoerceTo(o.Type()); ok {
			return o.Equal(cast)
		} else {
			return false, false
		}
	default:
		return false, false
	}
}

var _ conformance.Comparable = (*StringObject)(nil)

// MARK: ComparisionEvaluator

func (o *StringObject) PerformEqualityCheck(r Object) (bool, error) {
	i, err := CoerceTo(r, STRING)
	if err != nil {
		return false, err
	}

	return o.Value == i.(*StringObject).Value, nil
}

func (o *StringObject) PerformLessThanComparison(r Object) (bool, error) {
	i, err := CoerceTo(r, STRING)
	if err != nil {
		return false, err
	}

	return o.Value < i.(*StringObject).Value, nil
}

var _ ComparisionEvaluator = (*StringObject)(nil)
