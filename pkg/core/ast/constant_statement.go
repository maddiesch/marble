package ast

import (
	"strings"

	"github.com/maddiesch/marble/pkg/core/token"
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

func (s *ConstantStatement) MarshalJSON() ([]byte, error) {
	return marshalNode(s)
}

func (*ConstantStatement) Name() string {
	return "ConstantStatement"
}

func (*ConstantStatement) _statementNode() {}

var _ Statement = (*ConstantStatement)(nil)

func (*ConstantStatement) Mutable() bool {
	return false
}

func (s *ConstantStatement) Label() string {
	return s.Identifier.Value
}

func (s *ConstantStatement) AssignmentExpression() Expression {
	return s.Expression
}

func (*ConstantStatement) CurrentFrameOnly() bool {
	return true
}

func (*ConstantStatement) RequireUndefined() bool {
	return true
}

func (*ConstantStatement) RequireDefined() bool {
	return false
}

var _ AssignmentStatement = (*ConstantStatement)(nil)
