package object

import "github.com/maddiesch/marble/pkg/core/visitor"

const (
	VOID ObjectType = "OBJ_Void"
)

var _voidObjectInstance = &VoidObject{}

func NewVoid() *VoidObject {
	return _voidObjectInstance
}

type VoidObject struct {
	Value bool
}

func (*VoidObject) Type() ObjectType {
	return VOID
}

func (*VoidObject) DebugString() string {
	return "void"
}

func (*VoidObject) GoValue() any {
	return nil
}

func (o *VoidObject) CoerceTo(t ObjectType) (Object, bool) {
	return NewVoid(), true
}

func (o *VoidObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*VoidObject)(nil)
