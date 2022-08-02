package ast

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/token"
)

type NotExpression struct {
	Token      token.Token
	Expression Expression
}

func (e *NotExpression) SourceToken() token.Token {
	return e.Token
}

func (e *NotExpression) String() string {
	return fmt.Sprintf("NOT(%s)", e.Expression)
}

func (b *NotExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*NotExpression) Name() string {
	return "BooleanNegateExpression"
}

func (*NotExpression) _expressionNode() {}

var _ Expression = (*NotExpression)(nil)
