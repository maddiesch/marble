package ast

import (
	"strings"

	"github.com/maddiesch/marble/pkg/token"
)

type IfExpression struct {
	Token          token.Token
	Condition      Expression
	TrueStatement  *BlockStatement
	FalseStatement *BlockStatement
}

func (e *IfExpression) SourceToken() token.Token {
	return e.Token
}

func (e *IfExpression) String() string {
	var builder strings.Builder

	builder.WriteString("if " + e.Condition.String())

	builder.WriteString(" " + e.TrueStatement.String() + " ")

	if e.FalseStatement != nil {
		builder.WriteString(" else " + e.FalseStatement.String() + " ")
	}

	return builder.String()
}

func (b *IfExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*IfExpression) Name() string {
	return "IfExpression"
}

func (*IfExpression) _expressionNode() {}

var _ Expression = (*IfExpression)(nil)
