package test

import (
	"fmt"
	"strings"

	"github.com/maddiesch/marble/internal/collection"
	"github.com/maddiesch/marble/pkg/core/ast"
	"github.com/maddiesch/marble/pkg/core/evaluator"
	"github.com/maddiesch/marble/pkg/core/lexer"
	"github.com/maddiesch/marble/pkg/core/object"
	"github.com/maddiesch/marble/pkg/core/parser"
	"github.com/stretchr/testify/require"
)

func Eval(t TestingT, source ...string) object.Object {
	pro := CreateProgram(t, source...)

	b := evaluator.NewBinding()

	out, err := evaluator.Evaluate(b, pro)

	require.NoError(t, err)

	return out
}

func CreateProgram(t TestingT, source ...string) *ast.Program {
	lexers := collection.MapSliceI(source, func(i int, s string) lexer.Lexer {
		l, err := lexer.New(fmt.Sprintf("%s+%d", t.Name(), i), strings.NewReader(s))

		require.NoError(t, err, "Failed to create lexer from source")

		return l
	})

	lex := lexer.NewLexer(lexers...)

	p := parser.New(lex)

	pro := p.Run()

	require.NoError(t, p.Err(), "Failed to parse program input")

	return pro
}
