package object

import "github.com/maddiesch/marble/pkg/core/visitor"

func GetSubscriptValue(from, key Object) (Object, *ErrorObject) {
	visitor := &SubscriptVisitor{Key: key}
	from.Accept(visitor)
	return visitor.Result, visitor.Error
}

type SubscriptVisitor struct {
	visitor.Visitor[Object]

	Key    Object
	Result Object
	Error  *ErrorObject
}

const (
	SubscriptVisitorError = "SubscriptError"
)

func (v *SubscriptVisitor) Visit(object Object) {
	switch object := object.(type) {
	case *ArrayObject:
		v.Result, v.Error = v.visitArray(object)
	default:
		v.Error = NewErrorf(SubscriptVisitorError, "Unable to perform subscript lookup for object %s", object.Type())
	}
}

func (v *SubscriptVisitor) visitArray(array *ArrayObject) (Object, *ErrorObject) {
	index, err := CastObjectTo(v.Key, INTEGER)
	if err != nil {
		return nil, err
	}
	return array.Elements[int(index.(*IntegerObject).Value)], nil
}
