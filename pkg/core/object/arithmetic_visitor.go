package object

import (
	"github.com/maddiesch/marble/pkg/core/visitor"
)

func GetArithmeticResult(op rune, left, right Object) (Object, *ErrorObject) {
	visitor := &ArithmeticVisitor{
		Operator:  op,
		RightHand: right,
	}

	left.Accept(visitor)

	return visitor.Result, visitor.Error
}

type ArithmeticVisitor struct {
	visitor.Visitor[Object]

	Operator  rune
	RightHand Object
	Result    Object
	Error     *ErrorObject
}

const (
	ArithmeticVisitorError = "ArithmeticError"
)

func (v *ArithmeticVisitor) Visit(in Object) {
	right, err := CastObjectTo(v.RightHand, in.Type())
	if err != nil {
		v.Error = err
		return
	}

	switch left := in.(type) {
	case *IntegerObject:
		v.visitInteger(left, right.(*IntegerObject))
	case *FloatObject:
		v.visitFloat(left, right.(*FloatObject))
	case *StringObject:
		v.visitString(left, right.(*StringObject))
	case *ArrayObject:
		v.visitArray(left, right.(*ArrayObject))
	default:
		v.Error = NewErrorf(ArithmeticVisitorError, "Arithmetic operation not supported for left-hand object with type: %s", left.Type())
	}
}

func (v *ArithmeticVisitor) visitArray(left, right *ArrayObject) {
	switch v.Operator {
	case '+':
		v.Result = NewArray(append(left.Elements, right.Elements...))
	default:
		v.assignUnsupportedOperatorError(v.Operator, STRING)
	}
}

func (v *ArithmeticVisitor) visitString(left, right *StringObject) {
	switch v.Operator {
	case '+':
		v.Result = NewString(left.Value + right.Value)
	default:
		v.assignUnsupportedOperatorError(v.Operator, STRING)
	}
}

func (v *ArithmeticVisitor) visitFloat(left, right *FloatObject) {
	switch v.Operator {
	case '+':
		v.Result = NewFloat(left.Value + right.Value)
	case '-':
		v.Result = NewFloat(left.Value - right.Value)
	case '*':
		v.Result = NewFloat(left.Value * right.Value)
	case '/':
		v.Result = NewFloat(left.Value / right.Value)
	default:
		v.assignUnsupportedOperatorError(v.Operator, FLOAT)
	}
}

func (v *ArithmeticVisitor) visitInteger(left, right *IntegerObject) {
	switch v.Operator {
	case '+':
		v.Result = NewInteger(left.Value + right.Value)
	case '-':
		v.Result = NewInteger(left.Value - right.Value)
	case '*':
		v.Result = NewInteger(left.Value * right.Value)
	case '/':
		v.Result = NewInteger(left.Value / right.Value)
	default:
		v.assignUnsupportedOperatorError(v.Operator, INTEGER)
	}
}

func (v *ArithmeticVisitor) assignUnsupportedOperatorError(op rune, kind ObjectType) {
	v.Error = NewErrorf(ArithmeticVisitorError, "Operator %s is not supported for object type %s", string(op), kind)
}
