package parser_test

import (
	"strings"
	"testing"

	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/core/ast"
	"github.com/maddiesch/marble/pkg/core/lexer"
	"github.com/maddiesch/marble/pkg/core/parser"
	"github.com/stretchr/testify/assert"
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

func TestParserEdgeCases(t *testing.T) {
	t.Run("if-else nested in block", func(t *testing.T) {
		par := createParserFromSource(t, `do { if (true) {} else {} }`)
		par.Run()

		assert.NoError(t, par.Err())
	})
}