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

func Closure(parameters []string, body *ast.BlockStatement, b *binding.Binding[Object]) *ClosureLiteral {
	return &ClosureLiteral{ParameterList: parameters, Body: body, Binding: b}
}

type ClosureLiteral struct {
	ParameterList []string
	Body          *ast.BlockStatement
	Binding       *binding.Binding[Object]
}

func (*ClosureLiteral) Type() ObjectType {
	return CLOSURE
}

func (o *ClosureLiteral) DebugString() string {
	return fmt.Sprintf("Closure(%d)", o.Binding.ID())
}

func (*ClosureLiteral) GoValue() any {
	return nil
}

func (*ClosureLiteral) CoerceTo(t ObjectType) (Object, bool) {
	return NewVoid(), false
}

func (o *ClosureLiteral) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*ClosureLiteral)(nil)
