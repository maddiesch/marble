package ast

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/token"
)

type FloatExpression struct {
	Token token.Token
	Value float64
}

func (e *FloatExpression) SourceToken() token.Token {
	return e.Token
}

func (e *FloatExpression) String() string {
	return fmt.Sprintf("Float(%f)", e.Value)
}

func (b *FloatExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*FloatExpression) Name() string {
	return "FloatExpression"
}

func (*FloatExpression) _expressionNode() {}

var _ Expression = (*FloatExpression)(nil)
