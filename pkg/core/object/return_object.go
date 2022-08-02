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

func (o *ReturnObject) CoerceTo(t ObjectType) (Object, bool) {
	return o.Value.CoerceTo(t)
}

func (o *ReturnObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*ReturnObject)(nil)

// MARK: ComparisionEvaluator

func (o *ReturnObject) PerformEqualityCheck(r Object) (bool, error) {
	if v, ok := o.Value.(EqualityEvaluator); ok {
		return v.PerformEqualityCheck(r)
	}
	return false, ReturnError{
		Value:   o.Value,
		Message: "Unable to perform equality check, as the return value object does not implement the equality interface",
	}
}

func (o *ReturnObject) PerformLessThanComparison(r Object) (bool, error) {
	if v, ok := o.Value.(ComparisionEvaluator); ok {
		return v.PerformLessThanComparison(r)
	}
	return false, ReturnError{
		Value:   o.Value,
		Message: "Unable to perform equality check, as the return value object does not implement the comparison interface",
	}
}

var _ ComparisionEvaluator = (*ReturnObject)(nil)

type ReturnError struct {
	Value   Object
	Message string
}

func (e ReturnError) Error() string {
	return fmt.Sprintf("ReturnError: {%s} %s", e.Value.DebugString(), e.Message)
}
