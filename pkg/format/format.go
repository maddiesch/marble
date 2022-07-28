package format

import "github.com/maddiesch/marble/pkg/ast"

func Format(node ast.Node) (string, error) {
	if f, ok := node.(Formattable); ok {
		return f.Format()
	}

	panic("format.Format - not implemented")
}

type Formattable interface {
	Format() (string, error)
}
