package object

import (
	"github.com/maddiesch/marble/pkg/core/visitor"
)

// CastVisitor handles converting between Object types.
//
// Note: The Result is _NOT_ guaranteed to be a new object instance, nor is is
// it guaranteed to be the same object instance.
type CastVisitor struct {
	visitor.Visitor[Object]

	Target ObjectType // The ObjectType the visited object should be converted into.

	Result Object       // The result of the conversion or nil
	Error  *ErrorObject // The conversion error or nil
}

func (v *CastVisitor) Take() (Object, *ErrorObject) {
	return v.Result, v.Error
}

const (
	ObjectCastError = "CastError"
)

func (v *CastVisitor) Visit(object Object) {
	switch object := object.(type) {
	case *IntegerObject:
		v.fromInteger(object)
	case *BoolObject:
		v.fromBoolean(object)
	case *FloatObject:
		v.fromFloat(object)
	case *NullObject:
		v.fromNull(object)
	case *ReturnObject:
		v.fromReturn(object)
	default:
		v.fromUnhandled(object)
	}
}

func (v *CastVisitor) fromReturn(object *ReturnObject) {
	child := &CastVisitor{Target: v.Target}

	object.Value.Accept(child)

	v.Result = child.Result
	v.Error = child.Error
}

func (v *CastVisitor) fromNull(object *NullObject) {
	switch v.Target {
	case NULL:
		v.Result = NewNull()
	case BOOLEAN:
		v.Result = NewBool(false)
	default:
		v.Error = NewErrorf(ObjectCastError, "Unable to cast from nil into %s, target is not supported", v.Target)
	}
}

func (v *CastVisitor) fromUnhandled(object Object) {
	if object.Type() == v.Target {
		v.Result = object
	} else {
		v.Error = NewErrorf(ObjectCastError, "Unable to cast from %s to %s, visited object type not supported", object.Type(), v.Target)
	}
}

func (v *CastVisitor) fromFloat(object *FloatObject) {
	switch v.Target {
	case FLOAT:
		v.Result = NewFloat(object.Value)
	case INTEGER:
		v.Result = NewInteger(int64(object.Value))
	default:
		v.Error = NewErrorf(ObjectCastError, "Unable to cast from float into %s, target is not supported", v.Target)
	}
}

func (v *CastVisitor) fromInteger(integer *IntegerObject) {
	switch v.Target {
	case BOOLEAN:
		v.Result = NewBool(integer.Value > 0)
	case FLOAT:
		v.Result = NewFloat(float64(integer.Value))
	case INTEGER:
		v.Result = NewInteger(integer.Value)
	default:
		v.Error = NewErrorf(ObjectCastError, "Unable to cast from integer into %s, target is not supported", v.Target)
	}
}

func (v *CastVisitor) fromBoolean(boolean *BoolObject) {
	switch v.Target {
	case BOOLEAN:
		v.Result = NewBool(boolean.Value)
	case FLOAT:
		if boolean.Value == true {
			v.Result = NewFloat(1.0)
		} else {
			v.Result = NewFloat(0.0)
		}
	case INTEGER:
		if boolean.Value == true {
			v.Result = NewInteger(1)
		} else {
			v.Result = NewInteger(0)
		}
	default:
		v.Error = NewErrorf(ObjectCastError, "Unable to cast from boolean into %s, target is not supported", v.Target)
	}
}
