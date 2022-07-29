package object

import "fmt"

const (
	INTEGER ObjectType = "OBJ_Integer"
)

type Integer struct {
	Value int64
}

func (*Integer) Type() ObjectType {
	return INTEGER
}

func (o *Integer) Description() string {
	return fmt.Sprintf("%d", o.Value)
}

func (o *Integer) GoValue() any {
	return o.Value
}

func (o *Integer) Cast(t ObjectType) (Object, bool) {
	switch t {
	case INTEGER:
		return o, true
	case FLOAT:
		return &Float{Value: float64(o.Value)}, true
	case BOOLEAN:
		return &Boolean{Value: o.Value > 0}, true
	default:
		return &Void{}, false
	}
}

var _ Object = (*Integer)(nil)
