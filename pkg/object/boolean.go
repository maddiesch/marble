package object

import "fmt"

const (
	BOOLEAN ObjectType = "OBJ_Boolean"
)

func Bool(v bool) *Boolean {
	return &Boolean{Value: v}
}

type Boolean struct {
	Value bool
}

func (o *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (o *Boolean) Description() string {
	return fmt.Sprintf("Bool(%t)", o.Value)
}

func (o *Boolean) GoValue() any {
	return o.Value
}

func (o *Boolean) CoerceTo(t ObjectType) (Object, bool) {
	switch t {
	case BOOLEAN:
		return o, true
	case INTEGER:
		if o.Value == true {
			return Int(1), true
		} else {
			return Int(0), true
		}
	case FLOAT:
		if o.Value == true {
			return Float(1.0), true
		} else {
			return Float(0.0), true
		}
	default:
		return &Void{}, false
	}
}

var _ Object = (*Boolean)(nil)

// MARK: EqualityEvaluator

func (o *Boolean) PerformEqualityCheck(r Object) (bool, error) {
	i, err := CoerceTo(r, BOOLEAN)
	if err != nil {
		return false, err
	}

	return o.Value == i.(*Boolean).Value, nil
}

var _ EqualityEvaluator = (*Boolean)(nil)
