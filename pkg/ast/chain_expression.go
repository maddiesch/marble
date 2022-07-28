package ast

import (
	"strings"

	"github.com/maddiesch/marble/pkg/token"
)

type ChainExpression struct {
	Token token.Token
	From  Expression
	To    Expression
}

func (e *ChainExpression) SourceToken() token.Token {
	return e.Token
}

func (e *ChainExpression) String() string {
	var builder strings.Builder
	builder.WriteString("CHAIN(" + e.From.String() + ")->(" + e.To.String() + "))")
	return builder.String()
}

func (b *ChainExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*ChainExpression) Name() string {
	return "ChainExpression"
}

func (*ChainExpression) _expressionNode() {}

var _ Expression = (*ChainExpression)(nil)
