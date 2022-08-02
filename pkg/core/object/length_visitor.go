package object

import "github.com/maddiesch/marble/pkg/core/visitor"

func GetObjectLength(object Object) (int, *ErrorObject) {
	visitor := &LengthVisitor{}

	object.Accept(visitor)

	return visitor.Length, visitor.Error
}

type LengthVisitor struct {
	visitor.Visitor[Object]

	Length int
	Error  *ErrorObject
}

const (
	LengthVisitorError = "LengthError"
)

func (v *LengthVisitor) Visit(object Object) {
	switch object := object.(type) {
	case *StringObject:
		v.Length = len(object.Value)
	case *ArrayObject:
		v.Length = len(object.Elements)
	default:
		v.Error = NewErrorf(LengthVisitorError, "Length is not supported for object %s", object.Type())
	}
}
