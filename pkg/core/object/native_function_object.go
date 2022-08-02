package object

import (
	"fmt"
	"strings"

	"github.com/maddiesch/marble/internal/collection"
	"github.com/maddiesch/marble/pkg/core/binding"
	"github.com/maddiesch/marble/pkg/core/visitor"
)

const (
	NATIVE_FUNC ObjectType = "OBJ_NativeClosure"
)

func NewNativeFunction(count int, fn NativeFunc) *NativeFunctionObject {
	return &NativeFunctionObject{ArgumentCount: count, Body: fn}
}

type NativeFunc func(*binding.Binding[Object], []Object) (Object, error)

type NativeFunctionObject struct {
	ArgumentCount int
	Body          NativeFunc
}

func (*NativeFunctionObject) Type() ObjectType {
	return NATIVE_FUNC
}

func (o *NativeFunctionObject) DebugString() string {
	return fmt.Sprintf("NativeFunc(%s)", strings.Join(collection.MapSlice(make([]string, o.ArgumentCount), func(string) string { return "_" }), ", "))
}

func (*NativeFunctionObject) GoValue() any {
	return nil
}

func (*NativeFunctionObject) CoerceTo(t ObjectType) (Object, bool) {
	return NewVoid(), false
}

func (o *NativeFunctionObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*NativeFunctionObject)(nil)
