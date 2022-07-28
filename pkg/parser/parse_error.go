package parser

import (
	"fmt"
	"strings"

	"github.com/maddiesch/marble/pkg/token"
)

type ParseError struct {
	Children []error
}

func (e *ParseError) Count() int {
	return len(e.Children)
}

func (e *ParseError) Error() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("ParseError: (%d)", len(e.Children)))
	for _, child := range e.Children {
		builder.WriteString(fmt.Sprintf("\n  %s", child.Error()))
	}

	return builder.String()
}

type LiteralParseError struct {
	Token           token.Token
	Message         string
	UnderlyingError error
}

func (e LiteralParseError) Error() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("LiteralParseError: %s", e.Message))
	builder.WriteString(fmt.Sprintf("\n  UnderlyingError: %s", e.UnderlyingError))
	builder.WriteString(fmt.Sprintf("\n  Source: %s", e.Token.Location))

	return builder.String()
}

type ExpressionParseError struct {
	Token   token.Token
	Message string
}

func (e ExpressionParseError) Error() string {
	return fmt.Sprintf("ExpressionParseError: %s\n  Token: %s\n  Source: %s", e.Message, e.Token, e.Token.Location)
}
