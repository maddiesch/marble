package parser_test

import (
	"strings"

	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/ast"
	"github.com/maddiesch/marble/pkg/lexer"
	"github.com/maddiesch/marble/pkg/parser"
	"github.com/stretchr/testify/require"
)

func createParserFromSource(t test.TestingT, source string) *parser.Parser {
	lex, err := lexer.New(t.Name(), strings.NewReader(source))

	require.NoError(t, err, "Failed to create lexer from source")

	return parser.New(lex)
}

func createProgramFromSource(t test.TestingT, source string) *ast.Program {
	p := createParserFromSource(t, source)

	prog := p.Run()

	require.NoError(t, p.Err())

	return prog
}
