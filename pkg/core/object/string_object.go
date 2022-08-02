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

// MARK: ComparisionEvaluator

func (o *StringObject) PerformEqualityCheck(r Object) (bool, error) {
	str, err := CastObjectTo(r, STRING)
	if err != nil {
		panic(err) // TODO: Better error handling
	}

	return o.Value == str.(*StringObject).Value, nil
}

func (o *StringObject) PerformLessThanComparison(r Object) (bool, error) {
	str, err := CastObjectTo(r, STRING)
	if err != nil {
		panic(err) // TODO: Better error handling
	}

	return o.Value < str.(*StringObject).Value, nil
}

var _ ComparisionEvaluator = (*StringObject)(nil)

func (s *StringObject) Concat(r Object) (Object, error) {
	str, err := CastObjectTo(r, STRING)
	if err != nil {
		panic(err) // TODO: Better error handling
	}

	return NewString(s.Value + str.(*StringObject).Value), nil
}

var _ ConcatingEvaluator = (*StringObject)(nil)

func (s *StringObject) Len() int {
	return len(s.Value)
}

var _ LengthEvaluator = (*StringObject)(nil)
