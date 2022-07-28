package ast

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/token"
)

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (e *InfixExpression) SourceToken() token.Token {
	return e.Token
}

func (e *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left, e.Operator, e.Right)
}

func (b *InfixExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*InfixExpression) Name() string {
	return "InfixExpression"
}

func (*InfixExpression) _expressionNode() {}

var _ Expression = (*InfixExpression)(nil)
