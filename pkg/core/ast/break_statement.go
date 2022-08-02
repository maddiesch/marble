package ast

import (
	"github.com/maddiesch/marble/pkg/core/token"
)

type BreakStatement struct {
	Token token.Token
}

func (s *BreakStatement) SourceToken() token.Token {
	return s.Token
}

func (s *BreakStatement) String() string {
	return "BREAK()"
}

func (s *BreakStatement) MarshalJSON() ([]byte, error) {
	return marshalNode(s)
}

func (*BreakStatement) Name() string {
	return "BlockStatement"
}

func (*BreakStatement) _statementNode() {}

var _ Statement = (*BreakStatement)(nil)
