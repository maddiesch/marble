package object

import "github.com/maddiesch/marble/pkg/core/visitor"

func GetObjectEquality(lhs, rhs Object) (bool, *ErrorObject) {
	visitor := &EqualityVisitor{
		RightHand: rhs,
	}

	lhs.Accept(visitor)

	return visitor.Equal, visitor.Error
}

type EqualityVisitor struct {
	visitor.Visitor[Object]

	RightHand Object

	Equal bool
	Error *ErrorObject
}

const (
	GetObjectEqualityError = "GetObjectEqualityError"
)

func (v *EqualityVisitor) Visit(object Object) {
	right, err := CastObjectTo(v.RightHand, object.Type())
	if err != nil {
		v.Error = err
		return
	}

	switch object := object.(type) {
	case *IntegerObject:
		v.Equal = object.Value == right.(*IntegerObject).Value
	case *FloatObject:
		v.Equal = object.Value == right.(*FloatObject).Value
	case *StringObject:
		v.Equal = object.Value == right.(*StringObject).Value
	case *BoolObject:
		v.Equal = object.Value == right.(*BoolObject).Value
	default:
		v.Error = NewErrorf(GetObjectEqualityError, "Left hand side object does not support equality (%s)", object.Type())
	}
}

func (v *EqualityVisitor) visitInteger(left, right *IntegerObject) {
	panic("visitInteger")
}
