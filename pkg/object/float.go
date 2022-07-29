package object

import "fmt"

const (
	FLOAT ObjectType = "OBJ_Float"
)

type Float struct {
	Value float64
}

func (*Float) Type() ObjectType {
	return FLOAT
}

func (o *Float) Description() string {
	return fmt.Sprintf("%f", o.Value)
}

func (o Float) GoValue() any {
	return o.Value
}

func (o *Float) Cast(t ObjectType) (Object, bool) {
	switch t {
	case INTEGER:
		return &Integer{Value: int64(o.Value)}, true
	case FLOAT:
		return o, true
	default:
		return &Void{}, false
	}
}

var _ Object = (*Float)(nil)
