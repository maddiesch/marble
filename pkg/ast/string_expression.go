package ast

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/token"
)

type StringExpression struct {
	Token token.Token
	Value string
}

func (e *StringExpression) SourceToken() token.Token {
	return e.Token
}

func (e *StringExpression) String() string {
	return fmt.Sprintf("String(%s)", e.Value)
}

func (b *StringExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*StringExpression) Name() string {
	return "StringExpression"
}

func (*StringExpression) _expressionNode() {}

var _ Expression = (*StringExpression)(nil)
