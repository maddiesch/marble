package object

import "fmt"

type Coercible interface {
	CoerceTo(ObjectType) (Object, bool)
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
		result, err = CoerceTo(o, t)
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
