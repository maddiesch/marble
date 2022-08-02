package parser

import (
	"strings"

	"github.com/maddiesch/marble/pkg/core/token"
)

type InvalidTokenError struct {
	Token token.Token
}

func (e InvalidTokenError) Error() string {
	var builder strings.Builder

	builder.WriteString("InvalidToken: Parser encountered an invalid token\n")
	builder.WriteString("  Source: " + e.Token.Location.String())

	return builder.String()
}
