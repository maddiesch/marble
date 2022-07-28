package ast

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/token"
)

type ExpressionStatement struct {
	Token token.Token

	Expression Expression
}

func (s *ExpressionStatement) SourceToken() token.Token {
	return s.Token
}

func (s *ExpressionStatement) String() string {
	return fmt.Sprintf("EXPR: %s", s.Expression.String())
}

func (b *ExpressionStatement) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*ExpressionStatement) Name() string {
	return "ExpressionStatement"
}

func (*ExpressionStatement) _statementNode() {}

var _ Statement = (*ExpressionStatement)(nil)
