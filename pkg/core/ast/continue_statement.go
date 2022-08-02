package ast

import (
	"github.com/maddiesch/marble/pkg/core/token"
)

type ContinueStatement struct {
	Token token.Token
}

func (s *ContinueStatement) SourceToken() token.Token {
	return s.Token
}

func (s *ContinueStatement) String() string {
	return "CONTINUE()"
}

func (s *ContinueStatement) MarshalJSON() ([]byte, error) {
	return marshalNode(s)
}

func (*ContinueStatement) Name() string {
	return "ContinueStatement"
}

func (*ContinueStatement) _statementNode() {}

var _ Statement = (*ContinueStatement)(nil)
