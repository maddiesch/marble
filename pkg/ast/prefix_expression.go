package ast

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/token"
)

type PrefixExpression struct {
	Token      token.Token
	Operator   string
	Expression Expression
}

func (e *PrefixExpression) SourceToken() token.Token {
	return e.Token
}

func (e *PrefixExpression) String() string {
	return fmt.Sprintf("%s(%s)", e.Operator, e.Expression)
}

func (b *PrefixExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*PrefixExpression) Name() string {
	return "PrefixExpression"
}

func (*PrefixExpression) _expressionNode() {}

var _ Expression = (*PrefixExpression)(nil)
