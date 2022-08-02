package ast

import (
	"strings"

	"github.com/maddiesch/marble/pkg/core/token"
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

func (*LetStatement) Mutable() bool {
	return true
}

func (s *LetStatement) Label() string {
	return s.Identifier.Value
}

func (s *LetStatement) AssignmentExpression() Expression {
	return s.Expression
}

func (*LetStatement) CurrentFrameOnly() bool {
	return true
}

func (s *LetStatement) RequireUndefined() bool {
	return true
}

func (*LetStatement) RequireDefined() bool {
	return false
}

var _ AssignmentStatement = (*LetStatement)(nil)
