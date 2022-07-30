package ast

import (
	"strings"

	"github.com/maddiesch/marble/internal/collection"
	"github.com/maddiesch/marble/pkg/token"
)

type CallExpression struct {
	Token token.Token

	Function  Expression // Identifier or FunctionExpression
	Arguments []Expression
}

func (e *CallExpression) SourceToken() token.Token {
	return e.Token
}

func (e *CallExpression) String() string {
	var builder strings.Builder

	args := collection.MapSlice(e.Arguments, func(a Expression) string {
		return a.String()
	})

	builder.WriteString("Call(" + e.Function.String() + "(")
	builder.WriteString(strings.Join(args, ", ") + "))")

	return builder.String()
}

func (e *CallExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(e)
}

func (*CallExpression) Name() string {
	return "CallExpression"
}

func (*CallExpression) _expressionNode() {}

var _ Expression = (*CallExpression)(nil)
