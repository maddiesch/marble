package ast

import (
	"strings"

	"github.com/maddiesch/marble/internal/collection"
	"github.com/maddiesch/marble/pkg/core/token"
)

type FunctionExpression struct {
	Token          token.Token
	Parameters     []*IdentifierExpression
	BlockStatement *BlockStatement
}

func (e *FunctionExpression) SourceToken() token.Token {
	return e.Token
}

func (e *FunctionExpression) String() string {
	var builder strings.Builder

	params := collection.MapSlice(e.Parameters, func(e *IdentifierExpression) string {
		return e.String()
	})

	builder.WriteString("FUNC(" + strings.Join(params, ", ") + ")")
	builder.WriteString(" " + e.BlockStatement.String())

	return builder.String()
}

func (b *FunctionExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*FunctionExpression) Name() string {
	return "FunctionExpression"
}

func (*FunctionExpression) _expressionNode() {}

var _ Expression = (*FunctionExpression)(nil)
