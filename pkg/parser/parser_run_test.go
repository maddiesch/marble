package parser_test

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/ast"
	"github.com/maddiesch/marble/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserRun(t *testing.T) {
	t.Run("when given an invalid token", func(t *testing.T) {
		p := createParserFromSource(t, "`")

		p.Run()

		err := p.Err().(*parser.ParseError)

		require.Error(t, err)

		if assert.Equal(t, 1, len(err.Children)) {
			assert.ErrorAs(t, err.Children[0], &parser.InvalidTokenError{})
		}
	})

	t.Run("parse if statement", func(t *testing.T) {
		pro := createProgramFromSource(t, `if (foo > 1) { bar }`)

		spew.Dump(pro)

		// TODO: Write Test for If Statement
	})

	t.Run("parse if else statement", func(t *testing.T) {
		pro := createProgramFromSource(t, `if (foo > 1) { bar } else { baz }`)

		spew.Dump(pro)

		// TODO: Write Test for If-Else Statement
	})

	t.Run("operator precedence", func(t *testing.T) {
		tests := []struct {
			src      string
			expected string
		}{
			{"true", "Bool(true)"},
			{"false", "Bool(false)"},
			{"3 > 5 == false", "((Int(3) > Int(5)) = Bool(false))"},
			{"3 < 5 == true", "((Int(3) < Int(5)) = Bool(true))"},
			{"not foo", "!(ID(foo))"},
			{"!foo", "!(ID(foo))"},
			{"-a * b", "(-(ID(a)) * ID(b))"},
			{"5 < 5", "(Int(5) < Int(5))"},
			{"a + b * c + 3.5 / e - f", "(((ID(a) + (ID(b) * ID(c))) + (Float(3.500000) / ID(e))) - ID(f))"},
			{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((Int(3) + (Int(4) * Int(5))) = ((Int(3) * Int(1)) + (Int(4) * Int(5))))"},
			{"1 + (2 + 3) + 4", "((Int(1) + (Int(2) + Int(3))) + Int(4))"},
			{"(5 + 5) * 2", "((Int(5) + Int(5)) * Int(2))"},
			{"-(5 + 5)", "-((Int(5) + Int(5)))"},
			{"!(true == true)", "!((Bool(true) = Bool(true)))"},
		}

		for _, tc := range tests {
			pro := createProgramFromSource(t, tc.src)

			if assert.Equal(t, 1, len(pro.StatementList)) {
				assert.Equal(t, fmt.Sprintf("EXPR: %s", tc.expected), pro.StatementList[0].String(), "Source: %s", tc.src)
			}
		}
	})

	t.Run("parse infix operators", func(t *testing.T) {
		p := createProgramFromSource(t, "5 > foo_bar;")

		infix := test.RequireType(t, &ast.ExpressionStatement{}, p.StatementList[0]).Expression

		test.AssertInfixExpression(t, infix, int64(5), ">", "foo_bar")
	})

	t.Run("parse bang and negative prefix expression", func(t *testing.T) {
		p := createParserFromSource(t, `!foobar; -42; not foobar;`)

		prog := p.Run()

		require.NoError(t, p.Err())
		require.Equal(t, 3, len(prog.StatementList))

		bangExpr := test.RequireType(t, &ast.ExpressionStatement{}, prog.StatementList[0])
		negExpr := test.RequireType(t, &ast.ExpressionStatement{}, prog.StatementList[1])
		notExpr := test.RequireType(t, &ast.ExpressionStatement{}, prog.StatementList[2])

		assert.Equal(t, "!", test.RequireType(t, &ast.PrefixExpression{}, bangExpr.Expression).Operator)
		assert.Equal(t, "-", test.RequireType(t, &ast.PrefixExpression{}, negExpr.Expression).Operator)
		assert.Equal(t, "!", test.RequireType(t, &ast.PrefixExpression{}, notExpr.Expression).Operator)
	})

	t.Run("parse basic expressions", func(t *testing.T) {
		p := createParserFromSource(t, `foobar; 42; 3.14; "foo bar"; true; false;`)

		prog := p.Run()

		require.NoError(t, p.Err())
		require.Equal(t, 6, len(prog.StatementList))

		idExpr := test.RequireType(t, &ast.ExpressionStatement{}, prog.StatementList[0])
		intExpr := test.RequireType(t, &ast.ExpressionStatement{}, prog.StatementList[1])
		floatExpr := test.RequireType(t, &ast.ExpressionStatement{}, prog.StatementList[2])
		stringExpr := test.RequireType(t, &ast.ExpressionStatement{}, prog.StatementList[3])
		trueExpr := test.RequireType(t, &ast.ExpressionStatement{}, prog.StatementList[4])
		falseExpr := test.RequireType(t, &ast.ExpressionStatement{}, prog.StatementList[5])

		assert.Equal(t, "foobar", test.RequireType(t, &ast.IdentifierExpression{}, idExpr.Expression).Value)
		assert.Equal(t, int64(42), test.RequireType(t, &ast.IntegerExpression{}, intExpr.Expression).Value)
		assert.Equal(t, float64(3.14), test.RequireType(t, &ast.FloatExpression{}, floatExpr.Expression).Value)
		assert.Equal(t, "foo bar", test.RequireType(t, &ast.StringExpression{}, stringExpr.Expression).Value)
		assert.True(t, test.RequireType(t, &ast.BooleanExpression{}, trueExpr.Expression).Value)
		assert.False(t, test.RequireType(t, &ast.BooleanExpression{}, falseExpr.Expression).Value)
	})

	t.Run("parse return statement", func(t *testing.T) {
		p := createParserFromSource(t, `return 1; return 1.5; return "bar"; return foo();`)

		prog := p.Run()

		require.NoError(t, p.Err())

		require.Equal(t, 4, len(prog.StatementList))

		test.RequireType(t, &ast.ReturnStatement{}, prog.StatementList[0])
		test.RequireType(t, &ast.ReturnStatement{}, prog.StatementList[1])
		test.RequireType(t, &ast.ReturnStatement{}, prog.StatementList[2])
		test.RequireType(t, &ast.ReturnStatement{}, prog.StatementList[3])
	})

	t.Run("parse let statement", func(t *testing.T) {
		prog := createProgramFromSource(t, `let foo = 1; const bar = "baz";`)

		letStmt := test.RequireType(t, &ast.LetStatement{}, prog.StatementList[0])
		test.AssertIdentifier(t, "foo", letStmt.Identifier)

		constStmt := test.RequireType(t, &ast.ConstantStatement{}, prog.StatementList[1])
		test.AssertIdentifier(t, "bar", constStmt.Identifier)
	})

	t.Run("parse const statement", func(t *testing.T) {
		prog := createProgramFromSource(t, `const STDOUT = os::file_ptr(1);`)

		letStmt := test.RequireType(t, &ast.ConstantStatement{}, prog.StatementList[0])
		test.AssertIdentifier(t, "STDOUT", letStmt.Identifier)

		spew.Dump(prog)
	})

	t.Run("parse mutate statement", func(t *testing.T) {
		prog := createProgramFromSource(t, `mut foo = 1;`)

		letStmt := test.RequireType(t, &ast.MutateStatement{}, prog.StatementList[0])
		test.AssertIdentifier(t, "foo", letStmt.Identifier)
	})

	t.Run("parse function", func(t *testing.T) {
		t.Run("with parameters", func(t *testing.T) {
			prog := createProgramFromSource(t, `fn(a, b) { return a + b }`)

			spew.Dump(prog)

			// TODO: Write Tests
		})

		t.Run("without parameters", func(t *testing.T) {
			prog := createProgramFromSource(t, `fn() { return false }`)

			spew.Dump(prog)

			// TODO: Write Tests
		})
	})

	t.Run("parse call", func(t *testing.T) {
		t.Run("named function", func(t *testing.T) {
			prog := createProgramFromSource(t, `add(1, 2 * 3, 4 + 5)`)

			spew.Dump(prog)

			// TODO: Write Tests
		})

		t.Run("inline function", func(t *testing.T) {
			prog := createProgramFromSource(t, `fn(a, b, c) { return a + b + c }(1, 2 * 3, 4 + 5)`)

			spew.Dump(prog)

			// TODO: Write Tests
		})
	})

	t.Run("parse doublecolon", func(t *testing.T) {
		prog := createProgramFromSource(t, `std::format("{}", "foo")`)

		spew.Dump(prog)

		// TODO: Write Tests
	})

	t.Run("parse period", func(t *testing.T) {
		prog := createProgramFromSource(t, `foo.bar.baz()`)

		spew.Dump(prog)

		// TODO: Write Tests
	})

	t.Run("parse do statement", func(t *testing.T) {
		t.Run("top level", func(t *testing.T) {
			prog := createProgramFromSource(t, `do { let foo = "test"; }`)

			spew.Dump(prog)

			// TODO: Write Tests
		})

		t.Run("in function", func(t *testing.T) {
			prog := createProgramFromSource(t, `
				fn() {
					do { let foo = "test" }
					let bar = 1
				}
			`)

			spew.Dump(prog)

			// TODO: Write Tests
		})
	})

	t.Run("parse assignment after non-semicolon if", func(t *testing.T) {
		prog := createProgramFromSource(t, `
		if (false) {}

		const foo = 1;
		`)

		require.Equal(t, 2, len(prog.StatementList))

		test.RequireType(t, &ast.IfExpression{},
			test.RequireType(t, &ast.ExpressionStatement{}, prog.StatementList[0]).Expression,
		)

		constant := test.RequireType(t, &ast.ConstantStatement{}, prog.StatementList[1])
		test.AssertIdentifier(t, "foo", constant.Identifier)
		test.AssertLiteralValue(t, int64(1), constant.Expression)
	})

	t.Run("parse nil literal", func(t *testing.T) {
		prog := createProgramFromSource(t, "nil")

		require.Equal(t, 1, len(prog.StatementList))

		test.RequireType(t, &ast.NullExpression{},
			test.RequireType(t, &ast.ExpressionStatement{}, prog.StatementList[0]).Expression,
		)
	})
}
