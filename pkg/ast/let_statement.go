package ast

import (
	"strings"

	"github.com/maddiesch/marble/pkg/token"
)

type LetStatement struct {
	Token token.Token

	Identifier *IdentifierExpression
	Expression Expression
}

func (s *LetStatement) SourceToken() token.Token {
	return s.Token
}

func (s *LetStatement) String() string {
	var builder strings.Builder

	builder.WriteString("LET(" + s.Identifier.String() + " = ")

	if s.Expression != nil {
		builder.WriteString(s.Expression.String())
	} else {
		builder.WriteString("<void>")
	}

	builder.WriteString(")")

	return builder.String()
}

func (b *LetStatement) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*LetStatement) Name() string {
	return "LetStatement"
}

func (*LetStatement) _statementNode() {}

var _ Statement = (*LetStatement)(nil)
