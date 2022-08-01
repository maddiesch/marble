package native

import (
	"github.com/maddiesch/marble/pkg/binding"
	"github.com/maddiesch/marble/pkg/evaluator/runtime"
	"github.com/maddiesch/marble/pkg/object"
	"github.com/maddiesch/marble/pkg/version"
)

var Functions = map[string]*object.NativeFunctionObject{
	"_marble_version":   object.NativeFunction(0, _marbleVersion),
	"print_description": object.NativeFunction(1, _printDescription),
	"len":               object.NativeFunction(1, _len),
}

func Bind(b *binding.Binding[object.Object]) {
	for k, fn := range Functions {
		b.Set(k, fn, binding.F_PROTECTED|binding.F_CONST|binding.F_NATIVE)
	}
}

func _printDescription(b *binding.Binding[object.Object], args []object.Object) (object.Object, error) {
	// TODO: Fix print description
	// io.WriteString(b.Stderr(), args[0].Description()+"\n")

	return args[0], nil
}

func _marbleVersion(*binding.Binding[object.Object], []object.Object) (object.Object, error) {
	return object.String(version.Current), nil
}

func _len(_ *binding.Binding[object.Object], args []object.Object) (object.Object, error) {
	if o, ok := args[0].(object.LengthEvaluator); ok {
		return object.Int(o.Len()), nil
	} else {
		return nil, runtime.NewError(runtime.ArgumentError, "Given type does not conform to length", runtime.ErrorValue("Type", args[0].Type()))
	}
}
