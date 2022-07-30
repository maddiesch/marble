package native

import (
	"io"
	"os"

	"github.com/maddiesch/marble/pkg/object"
	"github.com/maddiesch/marble/pkg/version"
)

var Functions = map[string]*object.NativeFunctionObject{
	"_marble_version":   object.NativeFunction(0, _marbleVersion),
	"print_description": object.NativeFunction(1, _printDescription),
}

func Bind(b object.Binding) {
	for k, fn := range Functions {
		b.SetProtected(k, fn)
	}
}

func _printDescription(_ object.Binding, args []object.Object) (object.Object, error) {
	io.WriteString(os.Stdout, args[0].Description()+"\n")

	return args[0], nil
}

func _marbleVersion(object.Binding, []object.Object) (object.Object, error) {
	return object.String(version.Current), nil
}
