package ast

import (
	"strings"

	"github.com/maddiesch/marble/pkg/token"
)

type DeleteStatement struct {
	Token token.Token

	Identifier *IdentifierExpression
}

func (s *DeleteStatement) SourceToken() token.Token {
	return s.Token
}

func (s *DeleteStatement) String() string {
	var builder strings.Builder

	builder.WriteString("DELETE(" + s.Identifier.String() + ")")

	return builder.String()
}

func (b *DeleteStatement) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*DeleteStatement) Name() string {
	return "DeleteStatement"
}

func (*DeleteStatement) _statementNode() {}

var _ Statement = (*DeleteStatement)(nil)
