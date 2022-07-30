package object

import (
	"fmt"
)

const (
	RETURN ObjectType = "OBJ_Return"
)

func Return(o Object) *ReturnObject {
	return &ReturnObject{Value: o}
}

type ReturnObject struct {
	Value Object
}

func (*ReturnObject) Type() ObjectType {
	return RETURN
}

func (o *ReturnObject) Description() string {
	return fmt.Sprintf("Return(%s)", o.Value)
}

func (o *ReturnObject) GoValue() any {
	return o.Value.GoValue()
}

func (o *ReturnObject) CoerceTo(t ObjectType) (Object, bool) {
	return o.Value.CoerceTo(t)
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
	return fmt.Sprintf("ReturnError: {%s} %s", e.Value.Description(), e.Message)
}