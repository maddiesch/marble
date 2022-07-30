package ast

import (
	"strings"

	"github.com/maddiesch/marble/pkg/token"
)

type DefinedExpression struct {
	Token token.Token

	Identifier *IdentifierExpression
}

func (s *DefinedExpression) SourceToken() token.Token {
	return s.Token
}

func (s *DefinedExpression) String() string {
	var builder strings.Builder

	builder.WriteString("DEFINED(" + s.Identifier.String() + ")")

	return builder.String()
}

func (b *DefinedExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*DefinedExpression) Name() string {
	return "DefinedStatement"
}

func (*DefinedExpression) _expressionNode() {}

var _ Expression = (*DefinedExpression)(nil)
