package ast

import (
	"strings"

	"github.com/maddiesch/marble/internal/collection"
	"github.com/maddiesch/marble/pkg/core/token"
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

	statements := collection.MapSlice(s.StatementList, func(s Statement) string {
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
