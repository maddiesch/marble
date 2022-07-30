package native

import (
	"io"

	"github.com/maddiesch/marble/pkg/evaluator/runtime"
	"github.com/maddiesch/marble/pkg/object"
	"github.com/maddiesch/marble/pkg/version"
)

var Functions = map[string]*object.NativeFunctionObject{
	"_marble_version":   object.NativeFunction(0, _marbleVersion),
	"print_description": object.NativeFunction(1, _printDescription),
	"len":               object.NativeFunction(1, _len),
}

func Bind(b object.Binding) {
	for k, fn := range Functions {
		b.SetProtected(k, fn)
	}
}

func _printDescription(b object.Binding, args []object.Object) (object.Object, error) {
	io.WriteString(b.Stderr(), args[0].Description()+"\n")

	return args[0], nil
}

func _marbleVersion(object.Binding, []object.Object) (object.Object, error) {
	return object.String(version.Current), nil
}

func _len(_ object.Binding, args []object.Object) (object.Object, error) {
	if o, ok := args[0].(object.LengthEvaluator); ok {
		return object.Int(o.Len()), nil
	} else {
		return nil, runtime.NewError(runtime.ArgumentError, "Given type does not conform to length", runtime.ErrorValue("Type", args[0].Type()))
	}
}
