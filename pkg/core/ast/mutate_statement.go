package ast

import (
	"strings"

	"github.com/maddiesch/marble/pkg/core/token"
)

type MutateStatement struct {
	Token token.Token

	Identifier *IdentifierExpression
	Expression Expression
}

func (s *MutateStatement) SourceToken() token.Token {
	return s.Token
}

func (s *MutateStatement) String() string {
	var builder strings.Builder

	builder.WriteString("MUTATE(" + s.Identifier.String() + " = ")

	if s.Expression != nil {
		builder.WriteString(s.Expression.String())
	} else {
		builder.WriteString("<void>")
	}

	builder.WriteString(")")

	return builder.String()
}

func (b *MutateStatement) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*MutateStatement) Name() string {
	return "MutateStatement"
}

func (*MutateStatement) _statementNode() {}

var _ Statement = (*MutateStatement)(nil)

func (s *MutateStatement) Label() string {
	return s.Identifier.Value
}

func (*MutateStatement) Mutable() bool {
	return true
}

func (s *MutateStatement) AssignmentExpression() Expression {
	return s.Expression
}

func (*MutateStatement) CurrentFrameOnly() bool {
	return false
}

func (*MutateStatement) RequireUndefined() bool {
	return false
}

func (*MutateStatement) RequireDefined() bool {
	return true
}

var _ AssignmentStatement = (*MutateStatement)(nil)
