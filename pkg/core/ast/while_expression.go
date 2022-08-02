package ast

import (
	"strings"

	"github.com/maddiesch/marble/pkg/core/token"
)

type WhileExpression struct {
	Token token.Token

	Condition Expression
	Block     *BlockStatement
}

func (o *WhileExpression) SourceToken() token.Token {
	return o.Token
}

func (o *WhileExpression) String() string {
	var builder strings.Builder

	builder.WriteString("WHILE(" + o.Condition.String() + ", " + o.Block.String() + ")")

	return builder.String()
}

func (o *WhileExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(o)
}

func (*WhileExpression) Name() string {
	return "WhileExpression"
}

func (*WhileExpression) _expressionNode() {}

var _ Expression = (*WhileExpression)(nil)
