package format

import (
	"bytes"
	"fmt"
	"io"

	"github.com/maddiesch/marble/pkg/ast"
)

func Format(node ast.Node) (string, error) {
	f := &Formatter{}
	if err := f.Format(node); err != nil {
		return "", err
	}
	return f.String(), nil
}

type Formattable interface {
	Format(*Formatter) error
}

type Formatter struct {
	buffer bytes.Buffer
}

func (f *Formatter) Format(node ast.Node) error {
	if node, ok := node.(Formattable); ok {
		return node.Format(f)
	}
	return fmt.Errorf("node %v does not implement Formattable interface", node)
}

func (f *Formatter) Write(p []byte) (int, error) {
	return f.buffer.Write(p)
}

func (f *Formatter) String() string {
	return f.buffer.String()
}

var _ io.Writer = (*Formatter)(nil)
