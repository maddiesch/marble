package object

import "github.com/maddiesch/marble/pkg/core/visitor"

const (
	NULL ObjectType = "OBJ_Null"
)

var _nullObjectInstance = &NullObject{}

func NewNull() *NullObject {
	return _nullObjectInstance
}

type NullObject struct {
}

func (o *NullObject) Type() ObjectType {
	return NULL
}

func (o *NullObject) DebugString() string {
	return "Null()"
}

func (o *NullObject) GoValue() any {
	return nil
}

func (o *NullObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*NullObject)(nil)
