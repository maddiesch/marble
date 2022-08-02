package parser

import (
	"fmt"
	"strings"

	"github.com/maddiesch/marble/pkg/core/token"
)

type UnexpectedTokenError struct {
	Token       token.Token
	Expected    token.Kind
	ExpectedNot *token.Kind
	Alternate   *token.Kind
}

func (e UnexpectedTokenError) Error() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("UnexpectedTokenError: Expected %s, got %s", e.Expected, e.Token.Kind))
	builder.WriteString(fmt.Sprintf("\n  Source: %s", e.Token.Location))

	return builder.String()
}
