package ast

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/token"
)

type NullExpression struct {
	Token token.Token
}

func (e *NullExpression) SourceToken() token.Token {
	return e.Token
}

func (e *NullExpression) String() string {
	return fmt.Sprintf("NULL()")
}

func (b *NullExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*NullExpression) Name() string {
	return "NullExpression"
}

func (*NullExpression) _expressionNode() {}

var _ Expression = (*NullExpression)(nil)
