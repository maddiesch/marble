package object

import "fmt"

const (
	BOOLEAN ObjectType = "OBJ_Boolean"
)

type Boolean struct {
	Value bool
}

func (o *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (o *Boolean) Description() string {
	return fmt.Sprintf("%t", o.Value)
}

func (o *Boolean) GoValue() any {
	return o.Value
}

func (o *Boolean) Cast(t ObjectType) (Object, bool) {
	if t == o.Type() {
		return o, true
	}
	return &Void{}, false
}

var _ Object = (*Boolean)(nil)
