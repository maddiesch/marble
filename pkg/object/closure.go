package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/ast"
)

const (
	CLOSURE ObjectType = "OBJ_Closure"
)

func Closure(parameters []string, body *ast.BlockStatement, frameID uint64) *ClosureLiteral {
	return &ClosureLiteral{ParameterList: parameters, Body: body, FrameID: frameID}
}

type ClosureLiteral struct {
	ParameterList []string
	Body          *ast.BlockStatement
	FrameID       uint64
}

func (*ClosureLiteral) Type() ObjectType {
	return CLOSURE
}

func (o *ClosureLiteral) Description() string {
	return fmt.Sprintf("Closure(%d)", o.FrameID)
}

func (*ClosureLiteral) GoValue() any {
	return nil
}

func (*ClosureLiteral) CoerceTo(t ObjectType) (Object, bool) {
	return &Void{}, false
}

var _ Object = (*ClosureLiteral)(nil)
