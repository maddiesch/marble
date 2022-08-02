package object

// TODO: Delete once CastVisitor is complete

import "fmt"

type Coercible interface {
	CoerceTo(ObjectType) (Object, bool)
}

func CoerceToType[T Object](from Object, val T) error {
	out, err := CoerceTo(from, val.Type())
	if err != nil {
		return err
	}
	if outSafe, ok := out.(T); ok {
		CopyObject(outSafe, val)
	} else {
		panic("coercion succeeded, but the returned object was not the expected type")
	}

	return nil
}

func CoerceTo(obj Object, t ObjectType) (Object, error) {
	if o, ok := obj.CoerceTo(t); ok {
		return o, nil
	}
	return nil, &CoercionError{
		Source: obj,
		Target: t,
	}
}

func ChainCoerceTo(o Object, t ...ObjectType) (result Object, err error) {
	result = o
	for _, t := range t {
		result, err = CoerceTo(result, t)
		if err != nil {
			return nil, err
		}
	}
	return
}

type CoercionError struct {
	Source Object
	Target ObjectType
}

func (e CoercionError) Error() string {
	return fmt.Sprintf("CoercionError: cannot convert value of type '%s' to type '%s' in coercion", e.Source.Type(), e.Target)
}
