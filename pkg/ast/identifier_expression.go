package ast

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/token"
)

type IdentifierExpression struct {
	Token token.Token
	Value string
}

func (e *IdentifierExpression) SourceToken() token.Token {
	return e.Token
}

func (e *IdentifierExpression) String() string {
	return fmt.Sprintf("ID(%s)", e.Value)
}

func (b *IdentifierExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*IdentifierExpression) Name() string {
	return "IdentifierExpression"
}

func (*IdentifierExpression) _expressionNode() {}

var _ Expression = (*IdentifierExpression)(nil)
