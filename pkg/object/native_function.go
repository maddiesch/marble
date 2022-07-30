package object

import (
	"fmt"
	"strings"

	"github.com/maddiesch/marble/internal/collection"
)

const (
	NATIVE_FUNC ObjectType = "OBJ_NativeClosure"
)

func NativeFunction(count int, fn NativeFunc) *NativeFunctionObject {
	return &NativeFunctionObject{ArgumentCount: count, Body: fn}
}

type NativeFunc func(Binding, []Object) (Object, error)

type NativeFunctionObject struct {
	ArgumentCount int
	Body          NativeFunc
}

func (*NativeFunctionObject) Type() ObjectType {
	return NATIVE_FUNC
}

func (o *NativeFunctionObject) Description() string {
	return fmt.Sprintf("NativeFunc(%s)", strings.Join(collection.MapSlice(make([]string, o.ArgumentCount), func(string) string { return "_" }), ", "))
}

func (*NativeFunctionObject) GoValue() any {
	return nil
}

func (*NativeFunctionObject) CoerceTo(t ObjectType) (Object, bool) {
	return &Void{}, false
}

var _ Object = (*NativeFunctionObject)(nil)
