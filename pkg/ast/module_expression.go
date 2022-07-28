package ast

import (
	"strings"

	"github.com/maddiesch/marble/pkg/token"
)

type ModuleExpression struct {
	Token      token.Token
	Module     *IdentifierExpression
	Expression Expression
}

func (e *ModuleExpression) SourceToken() token.Token {
	return e.Token
}

func (e *ModuleExpression) String() string {
	var builder strings.Builder
	builder.WriteString("MOD(" + e.Module.String() + ")->(" + e.Expression.String() + "))")
	return builder.String()
}

func (b *ModuleExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*ModuleExpression) Name() string {
	return "ModuleExpression"
}

func (*ModuleExpression) _expressionNode() {}

var _ Expression = (*ModuleExpression)(nil)
