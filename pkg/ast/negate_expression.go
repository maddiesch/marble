package ast

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/token"
)

type NegateExpression struct {
	Token      token.Token
	Expression Expression
}

func (e *NegateExpression) SourceToken() token.Token {
	return e.Token
}

func (e *NegateExpression) String() string {
	return fmt.Sprintf("NEG(%s)", e.Expression)
}

func (b *NegateExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*NegateExpression) Name() string {
	return "NegateExpression"
}

func (*NegateExpression) _expressionNode() {}

var _ Expression = (*NegateExpression)(nil)
