package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/ast"
	"github.com/maddiesch/marble/pkg/core/binding"
	"github.com/maddiesch/marble/pkg/core/visitor"
)

const (
	CLOSURE ObjectType = "OBJ_Closure"
)

func NewClosure(parameters []string, body *ast.BlockStatement, b *binding.Binding[Object]) *ClosureObject {
	return &ClosureObject{ParameterList: parameters, Body: body, Binding: b}
}

type ClosureObject struct {
	ParameterList []string
	Body          *ast.BlockStatement
	Binding       *binding.Binding[Object]
}

func (*ClosureObject) Type() ObjectType {
	return CLOSURE
}

func (o *ClosureObject) DebugString() string {
	return fmt.Sprintf("Closure(%d)", o.Binding.ID())
}

func (*ClosureObject) GoValue() any {
	return nil
}

func (*ClosureObject) CoerceTo(t ObjectType) (Object, bool) {
	return NewVoid(), false
}

func (o *ClosureObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*ClosureObject)(nil)
