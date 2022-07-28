package ast

import (
	"strings"

	"github.com/maddiesch/marble/pkg/token"
)

type ConstantStatement struct {
	Token token.Token

	Identifier *IdentifierExpression
	Expression Expression
}

func (s *ConstantStatement) SourceToken() token.Token {
	return s.Token
}

func (s *ConstantStatement) String() string {
	var builder strings.Builder

	builder.WriteString("CONST(" + s.Identifier.String() + " = ")

	if s.Expression != nil {
		builder.WriteString(s.Expression.String())
	} else {
		builder.WriteString("<void>")
	}

	builder.WriteString(")")

	return builder.String()
}

func (b *ConstantStatement) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*ConstantStatement) Name() string {
	return "ConstantStatement"
}

func (*ConstantStatement) _statementNode() {}

var _ Statement = (*ConstantStatement)(nil)
