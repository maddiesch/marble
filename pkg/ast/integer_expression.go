package ast

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/token"
)

type IntegerExpression struct {
	Token token.Token
	Value int64
}

func (e *IntegerExpression) SourceToken() token.Token {
	return e.Token
}

func (e *IntegerExpression) String() string {
	return fmt.Sprintf("Int(%d)", e.Value)
}

func (b *IntegerExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*IntegerExpression) Name() string {
	return "IntegerExpression"
}

func (*IntegerExpression) _expressionNode() {}

var _ Expression = (*IntegerExpression)(nil)
