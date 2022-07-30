package ast

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/token"
)

type SubscriptExpression struct {
	Token    token.Token
	Receiver Expression
	Value    Expression
}

func (e *SubscriptExpression) SourceToken() token.Token {
	return e.Token
}

func (e *SubscriptExpression) String() string {
	return fmt.Sprintf("SUBSCRIPT(%s, %s)", e.Receiver, e.Value)
}

func (e *SubscriptExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(e)
}

func (*SubscriptExpression) Name() string {
	return "SubscriptExpression"
}

func (*SubscriptExpression) _expressionNode() {}

var _ Expression = (*SubscriptExpression)(nil)
