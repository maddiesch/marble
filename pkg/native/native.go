package native

import (
	"fmt"
	"io"
	"os"

	"github.com/maddiesch/marble/pkg/build"
	"github.com/maddiesch/marble/pkg/core/binding"
	"github.com/maddiesch/marble/pkg/core/evaluator/runtime"
	"github.com/maddiesch/marble/pkg/core/object"
)

var Functions = map[string]*object.NativeFunctionObject{
	"print_description": object.NativeFunction(1, _printDescription),
	"len":               object.NativeFunction(1, _len),
	"marble_dump":       object.NativeFunction(0, _marbleDump),
}

var Constants = map[string]object.Object{
	"MARBLE_VERSION": object.String(build.Version),
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
	if o, ok := args[0].(object.LengthEvaluator); ok {
		return object.Int(o.Len()), nil
	} else {
		return nil, runtime.NewError(runtime.ArgumentError, "Given type does not conform to length", runtime.ErrorValue("Type", args[0].Type()))
	}
}

func _marbleDump(b *binding.Binding[object.Object], _ []object.Object) (object.Object, error) {
	fmt.Fprintf(os.Stdout, b.DebugString())
	return &object.Void{}, nil
}
