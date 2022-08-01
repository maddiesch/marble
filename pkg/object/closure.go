package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/ast"
	"github.com/maddiesch/marble/pkg/binding"
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

func (o *ClosureLiteral) Description() string {
	return fmt.Sprintf("Closure(%d)", o.Binding.ID())
}

func (*ClosureLiteral) GoValue() any {
	return nil
}

func (*ClosureLiteral) CoerceTo(t ObjectType) (Object, bool) {
	return &Void{}, false
}

var _ Object = (*ClosureLiteral)(nil)
