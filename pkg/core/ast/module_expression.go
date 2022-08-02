package ast

import (
	"strings"

	"github.com/maddiesch/marble/pkg/core/token"
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
	builder.WriteString("MOD(")
	builder.WriteString("NAME(" + e.Module.String() + ")")
	builder.WriteString(", TARGET(" + e.Expression.String() + ")")
	builder.WriteString(")")
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
