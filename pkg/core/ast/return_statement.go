package ast

import (
	"strings"

	"github.com/maddiesch/marble/pkg/core/token"
)

type ReturnStatement struct {
	Token token.Token

	Expression Expression
}

func (s *ReturnStatement) SourceToken() token.Token {
	return s.Token
}

func (s *ReturnStatement) String() string {
	var builder strings.Builder

	builder.WriteString("RETURN(")
	if s.Expression != nil {
		builder.WriteString(s.Expression.String())
	} else {
		builder.WriteString("<void>")
	}
	builder.WriteString(")")

	return builder.String()
}

func (b *ReturnStatement) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*ReturnStatement) Name() string {
	return "ReturnStatement"
}

func (*ReturnStatement) _statementNode() {}

var _ Statement = (*ReturnStatement)(nil)
