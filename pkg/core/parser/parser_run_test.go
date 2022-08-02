package parser_test

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/core/ast"
	"github.com/maddiesch/marble/pkg/core/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserRun(t *testing.T) {
	t.Run("continue", func(t *testing.T) {
		pro := createProgramFromSource(t, `while (true) { continue }`)

		assert.Equal(t, `STMT(WHILE(Bool(true), { CONTINUE() }))`, pro.StatementList[0].String())
	})

	t.Run("break", func(t *testing.T) {
		pro := createProgramFromSource(t, `while (true) { break }`)

		assert.Equal(t, `STMT(WHILE(Bool(true), { BREAK() }))`, pro.StatementList[0].String())
	})

	t.Run("mutate assignment", func(t *testing.T) {
		pro := createProgramFromSource(t, `count = count + 1`)

		assert.Equal(t, `STMT(EXPR(MUTATE(ID(count) = (ID(count) + Int(1)))))`, pro.StatementList[0].String())
	})

	t.Run("while", func(t *testing.T) {
		pro := createProgramFromSource(t, `while (false) { let v = 0; }`)

		require.Equal(t, 1, len(pro.StatementList))

		expr := test.RequireType(t, &ast.ExpressionStatement{}, pro.StatementList[0]).Expression

		assert.Equal(t, "WHILE(Bool(false), { LET(ID(v) = Int(0)) })", expr.String())
	})

	t.Run("subscript", func(t *testing.T) {
		pro := createProgramFromSource(t, `value[1]`)

		require.Equal(t, 1, len(pro.StatementList))

		expr := test.RequireType(t, &ast.ExpressionStatement{}, pro.StatementList[0]).Expression

		assert.Equal(t, "SUBSCRIPT(ID(value), Int(1))", expr.String())
	})

	t.Run("array literal", func(t *testing.T) {
		pro := createProgramFromSource(t, `[1, 3.14, "Hello World", fn() { true }]`)

		require.Equal(t, 1, len(pro.StatementList))

		array := test.RequireType(t, &ast.ArrayExpression{},
			test.RequireType(t, &ast.ExpressionStatement{}, pro.StatementList[0]).Expression,
		)

		require.Equal(t, 4, len(array.Elements))

		assert.IsType(t, &ast.IntegerExpression{}, array.Elements[0])
		assert.IsType(t, &ast.FloatExpression{}, array.Elements[1])
		assert.IsType(t, &ast.StringExpression{}, array.Elements[2])
		assert.IsType(t, &ast.FunctionExpression{}, array.Elements[3])
	})

	t.Run("defined statement", func(t *testing.T) {
		t.Run("top level expression", func(t *testing.T) {
			pro := createProgramFromSource(t, `defined foo;`)

			require.Equal(t, 1, len(pro.StatementList))

			assert.Equal(t, `STMT(DEFINED(ID(foo)))`, pro.StatementList[0].String())
		})
	})

	t.Run("delete statement", func(t *testing.T) {
		pro := createProgramFromSource(t, `let foo = 1; delete foo; return true;`)

		require.Equal(t, 3, len(pro.StatementList))

		assert.Equal(t, `DELETE(ID(foo))`, pro.StatementList[1].String())
	})

	t.Run("when given an invalid token", func(t *testing.T) {
		p := createParserFromSource(t, "`")

		p.Run()

		err := p.Err().(*parser.ParseError)

		require.Error(t, err)

		if assert.Equal(t, 1, len(err.Children)) {
			assert.ErrorAs(t, err.Children[0], &parser.InvalidTokenError{})
		}
	})

	t.Run("operator precedence", func(t *testing.T) {
		tests := test.TestingTuple2[string, string]{
			{One: "true", Two: "Bool(true)"},
			{One: "false", Two: "Bool(false)"},
			{One: "3 > 5 == false", Two: "((Int(3) > Int(5)) == Bool(false))"},
			{One: "3 < 5 == true", Two: "((Int(3) < Int(5)) == Bool(true))"},
			{One: "not foo", Two: "NOT(ID(foo))"},
			{One: "!foo", Two: "NOT(ID(foo))"},
			{One: "-a * b", Two: "(NEG(ID(a)) * ID(b))"},
			{One: "5 < 5", Two: "(Int(5) < Int(5))"},
			{One: "a + b * c + 3.5 / e - f", Two: "(((ID(a) + (ID(b) * ID(c))) + (Float(3.500000) / ID(e))) - ID(f))"},
			{One: "3 + 4 * 5 == 3 * 1 + 4 * 5", Two: "((Int(3) + (Int(4) * Int(5))) == ((Int(3) * Int(1)) + (Int(4) * Int(5))))"},
			{One: "1 + (2 + 3) + 4", Two: "((Int(1) + (Int(2) + Int(3))) + Int(4))"},
			{One: "(5 + 5) * 2", Two: "((Int(5) + Int(5)) * Int(2))"},
			{One: "-(5 + 5)", Two: "NEG((Int(5) + Int(5)))"},
			{One: "!(true == true)", Two: "NOT((Bool(true) == Bool(true)))"},
			{One: "a * [1, 2, 3, 4][b * c] * d", Two: "((ID(a) * SUBSCRIPT(ARRAY[Int(1), Int(2), Int(3), Int(4)], (ID(b) * ID(c)))) * ID(d))"},
		}

		tests.Each(func(source, expected string) {
			t.Run(source, func(t *testing.T) {
				pro := createProgramFromSource(t, source)

				if assert.Equal(t, 1, len(pro.StatementList)) {
					assert.Equal(t, fmt.Sprintf("STMT(%s)", expected), pro.StatementList[0].String(), "Source: %s", source)
				}
			})
		})
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

		test.RequireType(t, &ast.NotExpression{}, bangExpr.Expression)
		test.RequireType(t, &ast.NegateExpression{}, negExpr.Expression)
		test.RequireType(t, &ast.NotExpression{}, notExpr.Expression)
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
