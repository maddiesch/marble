package native

import (
	"fmt"
	"io"
	"os"

	"github.com/maddiesch/marble/pkg/core/binding"
	"github.com/maddiesch/marble/pkg/core/object"
	"github.com/maddiesch/marble/pkg/version"
)

var Functions = map[string]*object.NativeFunctionObject{
	"print_description": object.NewNativeFunction(1, _printDescription),
	"len":               object.NewNativeFunction(1, _len),
	"marble_dump":       object.NewNativeFunction(0, _marbleDump),
}

var Constants = map[string]object.Object{
	"MARBLE_VERSION": object.NewString(version.Current),
}

func Bind(b *binding.Binding[object.Object]) {
	for k, fn := range Functions {
		b.Set(k, fn, binding.F_PROTECTED|binding.F_CONST|binding.F_NATIVE)
	}
	for k, v := range Constants {
		b.Set(k, v, binding.F_PROTECTED|binding.F_CONST|binding.F_NATIVE)
	}
}

func _printDescription(b *binding.Binding[object.Object], args []object.Object) (object.Object, error) {
	io.WriteString(os.Stdout, args[0].DebugString()+"\n")

	return args[0], nil
}

func _len(_ *binding.Binding[object.Object], args []object.Object) (object.Object, error) {
	if len, err := object.GetObjectLength(args[0]); err != nil {
		panic(err) // TODO: Better error handling
	} else {
		return object.NewInteger(len), nil
	}
}

func _marbleDump(b *binding.Binding[object.Object], _ []object.Object) (object.Object, error) {
	fmt.Fprintf(os.Stdout, b.DebugString())
	return object.NewVoid(), nil
}
