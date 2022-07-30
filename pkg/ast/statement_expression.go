package ast

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/token"
)

type StatementExpression struct {
	Token token.Token

	Statement Statement
}

func (s *StatementExpression) SourceToken() token.Token {
	return s.Token
}

func (s *StatementExpression) String() string {
	return fmt.Sprintf("EXPR(%s)", s.Statement.String())
}

func (b *StatementExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*StatementExpression) Name() string {
	return "StatementExpression"
}

func (*StatementExpression) _expressionNode() {}

var _ Expression = (*StatementExpression)(nil)
