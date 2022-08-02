package object

import "github.com/maddiesch/marble/pkg/core/visitor"

// GetLessThanComparison performs a less-than comparison between the to given objects
//
// It uses the left hand side as the casting target for the right hand side
func GetLessThanComparison(lhs, rhs Object) (bool, *ErrorObject) {
	visitor := &ComparisonVisitor{RightHand: rhs}
	lhs.Accept(visitor)
	return visitor.LessThan, visitor.Error
}

type ComparisonVisitor struct {
	visitor.Visitor[Object]

	RightHand Object

	LessThan bool
	Error    *ErrorObject
}

const (
	ComparisonVisitorError = "ComparisonError"
)

func (v *ComparisonVisitor) Visit(object Object) {
	switch object := object.(type) {
	case *IntegerObject:
		v.visitInteger(object)
	case *FloatObject:
		v.visitFloat(object)
	case *StringObject:
		v.visitString(object)
	default:
		v.Error = NewErrorf(ComparisonVisitorError, "Unable to perform less-than comparison (%s < %s)", object.Type(), v.RightHand.Type())
	}
}

func (v *ComparisonVisitor) visitString(object *StringObject) {
	if right, err := CastObjectTo(v.RightHand, STRING); err == nil {
		v.LessThan = object.Value < right.(*StringObject).Value
	} else {
		v.Error = err
	}
}

func (v *ComparisonVisitor) visitFloat(object *FloatObject) {
	if right, err := CastObjectTo(v.RightHand, FLOAT); err == nil {
		v.LessThan = object.Value < right.(*FloatObject).Value
	} else {
		v.Error = err
	}
}

func (v *ComparisonVisitor) visitInteger(object *IntegerObject) {
	if right, err := CastObjectTo(v.RightHand, INTEGER); err == nil {
		v.LessThan = object.Value < right.(*IntegerObject).Value
	} else {
		v.Error = err
	}
}
