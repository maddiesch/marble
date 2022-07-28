package ast

import (
	"strings"

	"github.com/maddiesch/marble/internal/slice"
	"github.com/maddiesch/marble/pkg/token"
)

type BlockStatement struct {
	Token token.Token

	StatementList []Statement
}

func (s *BlockStatement) SourceToken() token.Token {
	return s.Token
}

func (s *BlockStatement) String() string {
	var builder strings.Builder

	builder.WriteString("{ ")

	statements := slice.Map(s.StatementList, func(s Statement) string {
		return s.String()
	})

	builder.WriteString(strings.Join(statements, "; "))

	builder.WriteString(" }")

	return builder.String()
}

func (b *BlockStatement) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*BlockStatement) Name() string {
	return "BlockStatement"
}

func (*BlockStatement) _statementNode() {}

var _ Statement = (*BlockStatement)(nil)
