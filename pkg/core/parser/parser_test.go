package parser_test

import (
	"strings"
	"testing"

	"github.com/maddiesch/marble/internal/collection"
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

func TestSyntax(t *testing.T) {
	type SyntaxTest struct {
		Name   string
		Input  string
		Output string
	}

	tests := []SyntaxTest{
		{
			Name:   "Let Variable Assignment",
			Input:  "let value = 1;",
			Output: "LET(ID(value) = Int(1))",
		},
		{
			Name:   "Const Variable Assignment",
			Input:  "const value = 1;",
			Output: "CONST(ID(value) = Int(1))",
		},
		{
			Name:   "Mutate Keyword Variable Assignment",
			Input:  "mutate value = 1;",
			Output: "MUTATE(ID(value) = Int(1))",
		},
		{
			Name:   "Mut Keyword Variable Assignment",
			Input:  "mut value = 1;",
			Output: "MUTATE(ID(value) = Int(1))",
		},
		{
			Name:   "Mutate Variable Assignment",
			Input:  "value = 1;",
			Output: "STMT(EXPR(MUTATE(ID(value) = Int(1))))",
		},
		{
			Name:   "Closure Assignment",
			Input:  "const print = fn(val) { print_debug(val) }",
			Output: "CONST(ID(print) = FUNC(ID(val)) { STMT(CALL(ID(print_debug)(ID(val)))) })",
		},
		{
			Name:   "If Statement",
			Input:  "if (foo > 1) { bar }",
			Output: "STMT(IF(COND((ID(foo) > Int(1))), TRUE({ STMT(ID(bar)) })))",
		},
		{
			Name:   "If Else Statement",
			Input:  "if (foo > 1) { bar } else { baz }",
			Output: "STMT(IF(COND((ID(foo) > Int(1))), TRUE({ STMT(ID(bar)) }), FALSE({ STMT(ID(baz)) })))",
		},
		{
			Name:   "Closure without parameters",
			Input:  "fn() { return nil }",
			Output: "STMT(FUNC() { RETURN(NULL()) })",
		},
		{
			Name:   "Closure with parameters",
			Input:  "fn(a, b) { return a + b }",
			Output: "STMT(FUNC(ID(a), ID(b)) { RETURN((ID(a) + ID(b))) })",
		},
		{
			Name:   "Return Defined",
			Input:  "return defined foo;",
			Output: "RETURN(DEFINED(ID(foo)))",
		},
		{
			Name:   "Call Named Function",
			Input:  "add(1, 2 * 3, 4 + 5)",
			Output: "STMT(CALL(ID(add)(Int(1), (Int(2) * Int(3)), (Int(4) + Int(5)))))",
		},
		{
			Name:   "Call Inline Function",
			Input:  "fn(a, b, c) { return a + b + c }(1, 2 * 3, 4 + 5)",
			Output: "STMT(CALL(FUNC(ID(a), ID(b), ID(c)) { RETURN(((ID(a) + ID(b)) + ID(c))) }(Int(1), (Int(2) * Int(3)), (Int(4) + Int(5)))))",
		},
		{
			Name:   "If Else in Do Block",
			Input:  "do { if (true) {} else {} }",
			Output: "STMT(DO { STMT(IF(COND(Bool(true)), TRUE({  }), FALSE({  }))) })",
		},
		{
			Name:   "Parse Module Expression",
			Input:  `std::format("{}", "foo")`,
			Output: "STMT(MOD(NAME(ID(std)), TARGET(CALL(ID(format)(String({}), String(foo))))))",
		},
		{
			Name:   "Chain Navigation",
			Input:  "foo.bar()",
			Output: "STMT(CHAIN(FROM(ID(foo)), TO(CALL(ID(bar)()))))",
		},
		{
			Name:   "Do Statement",
			Input:  `do { let foo = "test"; }`,
			Output: "STMT(DO { LET(ID(foo) = String(test)) })",
		},
		{
			Name:   "Do Statement in Closure",
			Input:  `fn() { do { let foo = "test" }; let bar = 1 }`,
			Output: "STMT(FUNC() { STMT(DO { LET(ID(foo) = String(test)) }); LET(ID(bar) = Int(1)) })",
		},
	}
	/**

	 */

	for _, c := range tests {
		t.Run(c.Name, func(t *testing.T) {
			lex, err := lexer.New(c.Name, strings.NewReader(c.Input))
			require.NoError(t, err, "failed to create source lexer")

			parse := parser.New(lex)
			program := parse.Run()
			require.NoError(t, parse.Err(), "failed to parse source program without error")

			statements := collection.MapSlice(program.StatementList, func(n ast.Statement) string {
				return n.String()
			})

			assert.Equal(t, c.Output, strings.Join(statements, "\n"))
		})
	}
}
