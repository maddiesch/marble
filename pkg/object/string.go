package object

import (
	"fmt"
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

func (s *StringObject) Concat(r Object) (Object, error) {
	r, err := CoerceTo(r, STRING)
	if err != nil {
		return nil, err
	}

	return String(s.Value + r.(*StringObject).Value), nil
}

var _ ConcatingEvaluator = (*StringObject)(nil)

func (s *StringObject) Len() int {
	return len(s.Value)
}

var _ LengthEvaluator = (*StringObject)(nil)
