package ast

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/token"
)

type BooleanExpression struct {
	Token token.Token
	Value bool
}

func (e *BooleanExpression) SourceToken() token.Token {
	return e.Token
}

func (e *BooleanExpression) String() string {
	return fmt.Sprintf("Bool(%t)", e.Value)
}

func (b *BooleanExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*BooleanExpression) Name() string {
	return "BooleanExpression"
}

func (*BooleanExpression) _expressionNode() {}

var _ Expression = (*BooleanExpression)(nil)
