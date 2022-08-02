package ast

import (
	"strings"

	"github.com/maddiesch/marble/internal/collection"
	"github.com/maddiesch/marble/pkg/core/token"
)

type ArrayExpression struct {
	Token    token.Token
	Elements []Expression
}

func (e *ArrayExpression) SourceToken() token.Token {
	return e.Token
}

func (e *ArrayExpression) String() string {
	var builder strings.Builder

	elements := collection.MapSlice(e.Elements, func(e Expression) string {
		return e.String()
	})

	builder.WriteString("ARRAY[" + strings.Join(elements, ", ") + "]")

	return builder.String()
}

func (e *ArrayExpression) MarshalJSON() ([]byte, error) {
	return marshalNode(e)
}

func (*ArrayExpression) Name() string {
	return "ArrayExpression"
}

func (*ArrayExpression) _expressionNode() {}

var _ Expression = (*ArrayExpression)(nil)
