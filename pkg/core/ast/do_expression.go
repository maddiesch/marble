package ast

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/token"
)

type DoExpression struct {
	Token token.Token
	Block *BlockStatement
}

func (e *DoExpression) SourceToken() token.Token {
	return e.Token
}

func (e *DoExpression) String() string {
	return fmt.Sprintf("DO %s", e.Block)
}

func (b *DoExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*DoExpression) Name() string {
	return "DoExpression"
}

func (*DoExpression) _expressionNode() {}

var _ Expression = (*DoExpression)(nil)
