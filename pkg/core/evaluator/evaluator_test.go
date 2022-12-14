package evaluator_test

import (
	"testing"
	"time"

	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/core/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	t.Run("All", func(t *testing.T) {
		tests := test.TestingTuple2[string, any]{
			{One: `let count = 1; count = count + 1; count`, Two: int64(2)},
			{One: `let array = [1, 2, "Hello World!"]; len(array);`, Two: int64(3)},
			{One: `let message = "Hello World!"; len(message);`, Two: int64(12)},
		}

		tests.Each(func(source string, expected any) {
			t.Run(source, func(t *testing.T) {
				result := test.Eval(t, source)

				assert.Equal(t, expected, result.GoValue())
			})
		})
	})

	t.Run("While", func(t *testing.T) {
		tests := test.TestingTuple2[string, any]{
			{One: `let count = 0; while (count <= 2) { count = count + 1; }; count`, Two: int64(3)},
			{One: `while (true == nil) { }`, Two: nil},
			{One: `while (1) { return 2; }`, Two: int64(2)},
			{One: `while (true) { if (true) { break } }; 1`, Two: int64(1)},
			{One: `let count = 0; while (count < 5) { count = count + 1; };`, Two: int64(5)},
			{One: `let count = 0; let value = 0; while (count < 2) { count = count + 1; if (count == 1) { continue }; value = value + 1 }; value`, Two: int64(1)},
		}

		tests.Each(func(source string, expected any) {
			t.Run(source, func(t *testing.T) {
				done := make(chan struct{})
				var result object.Object

				go func() {
					defer close(done)

					result = test.Eval(t, source)
				}()

				select {
				case <-time.After(500 * time.Millisecond):
					require.Fail(t, "Test took too long to execute! Never broke from while loop?")
				case <-done:
				}

				assert.Equal(t, expected, result.GoValue())
			})
		})
	})

	t.Run("Array", func(t *testing.T) {
		tests := test.TestingTuple2[string, any]{
			{One: `let a = [1, fn() { true }]; a[1]()`, Two: true},
		}

		tests.Each(func(source string, expected any) {
			t.Run(source, func(t *testing.T) {
				result := test.Eval(t, source)

				assert.Equal(t, expected, result.GoValue())
			})
		})
	})

	t.Run("String", func(t *testing.T) {
		tests := test.TestingTuple2[string, any]{
			{One: `"Hello" + " " + "World" + "!"`, Two: "Hello World!"},
			{One: `len("Hello World!")`, Two: int64(12)},
		}

		tests.Each(func(source string, expected any) {
			t.Run(source, func(t *testing.T) {
				result := test.Eval(t, source)

				assert.Equal(t, expected, result.GoValue())
			})
		})
	})

	t.Run("NativeFunction", func(t *testing.T) {
		tests := test.TestingTuple2[string, any]{
			{One: `print_description("Hello World!")`, Two: "Hello World!"},
		}

		tests.Each(func(source string, expected any) {
			t.Run(source, func(t *testing.T) {
				result := test.Eval(t, source)

				assert.Equal(t, expected, result.GoValue())
			})
		})
	})

	t.Run("Function", func(t *testing.T) {
		tests := test.TestingTuple2[string, any]{
			{One: `const add = fn(a, b) { a + b }; add(1, 5)`, Two: int64(6)},
			{One: `const add5 = fn(a) { a() + 5 }; add5(fn() { 10 });`, Two: int64(15)},
			{One: `fn(n) {}(5)`, Two: nil},
			{One: `let factorial = fn(n) { if (n == 0) { 1 } else { n * factorial(n - 1) } }; factorial(5)`, Two: int64(120)},
		}

		tests.Each(func(source string, expected any) {
			t.Run(source, func(t *testing.T) {
				result := test.Eval(t, source)

				assert.Equal(t, expected, result.GoValue())
			})
		})
	})

	t.Run("Defined", func(t *testing.T) {
		tests := test.TestingTuple2[string, any]{
			{One: "let a = 1; defined a;", Two: true},
			{One: "defined a;", Two: false},
		}

		tests.Each(func(source string, expected any) {
			t.Run(source, func(t *testing.T) {
				result := test.Eval(t, source)

				assert.Equal(t, expected, result.GoValue())
			})
		})
	})

	t.Run("Delete", func(t *testing.T) {
		tests := test.TestingTuple2[string, any]{
			{One: "let a = 1; return defined a;", Two: true},
			{One: "let a = 1; delete a; return defined a;", Two: false},
		}

		tests.Each(func(source string, expected any) {
			t.Run(source, func(t *testing.T) {
				result := test.Eval(t, source)

				assert.Equal(t, expected, result.GoValue())
			})
		})
	})

	t.Run("Assignment", func(t *testing.T) {
		tests := test.TestingTuple2[string, any]{
			{One: "let a = 1; return a;", Two: int64(1)},
			{One: "const a = 42; let b = 8; return a;", Two: int64(42)},
			{One: "let a = 1; do { a = a + 1 }; a", Two: int64(2)},
		}

		tests.Each(func(source string, expected any) {
			t.Run(source, func(t *testing.T) {
				result := test.Eval(t, source)

				assert.Equal(t, expected, result.GoValue())
			})
		})
	})

	t.Run("Return", func(t *testing.T) {
		tests := test.TestingTuple2[string, any]{
			{One: "1; return 2; 3;", Two: int64(2)},
			{One: "return (1 + 4);", Two: int64(5)},
			{One: "return (1 + 4) == true;", Two: false},
			{One: "if (10 > 1) { if (10 > 1) { return 10 } return 1 }", Two: int64(10)},
		}

		tests.Each(func(source string, expected any) {
			t.Run(source, func(t *testing.T) {
				result := test.Eval(t, source)

				assert.Equal(t, expected, result.GoValue())
			})
		})
	})

	t.Run("Conditionals", func(t *testing.T) {
		tests := test.TestingTuple2[string, any]{
			{One: "if (true) { 10 } else { 42 }", Two: int64(10)},
			{One: "if (false) { 10 } else { 42 }", Two: int64(42)},
			{One: "if (false) { 10 }", Two: nil},
			{One: `if ("test" == "test") { 10 }`, Two: int64(10)},
		}

		tests.Each(func(source string, expected any) {
			t.Run(source, func(t *testing.T) {
				result := test.Eval(t, source)

				assert.Equal(t, expected, result.GoValue())
			})
		})
	})

	t.Run("Integer Expression", func(t *testing.T) {
		t.Run("boolean results", func(t *testing.T) {
			tests := []test.Tuple2[bool, string]{
				{One: true, Two: "50 > 25"},
				{One: true, Two: "42 > 16"},
				{One: true, Two: "42 >= 42"},
				{One: true, Two: "42 < 100"},
				{One: true, Two: "42 <= 100.8"},
				{One: true, Two: "42 <= 42"},
				{One: true, Two: "42 == 42"},
				{One: false, Two: "42 != 42"},
			}

			for _, tu := range tests {
				t.Run(tu.Two, func(t *testing.T) {
					out := test.Eval(t, tu.Two)

					assert.Equal(t, object.BOOLEAN, out.Type())
					assert.Equal(t, tu.One, out.GoValue())
				})
			}
		})

		t.Run("integer results", func(t *testing.T) {
			tests := []test.Tuple2[int64, string]{
				{One: 42, Two: "42"},
				{One: -128, Two: "-128"},
				{One: 100, Two: "50 + 50"},
				{One: 42, Two: "100 - 58.0"},
				{One: 30, Two: "5 * 6"},
				{One: -30, Two: "5 * -6"},
				{One: -30, Two: "5 * -6"},
				{One: 52, Two: "(8 + 5) * 4"},
				{One: 26, Two: "((8 + 5) * 4) / 2"},
				{One: 156, Two: "6 * (((8 + 5) * 4) / 2)"},
			}

			for _, tu := range tests {
				out := test.Eval(t, tu.Two)

				assert.Equal(t, object.INTEGER, out.Type())
				assert.Equal(t, tu.One, out.GoValue())
			}
		})
	})

	t.Run("FloatExpression", func(t *testing.T) {
		tests := []test.Tuple2[float64, string]{
			{One: 42.0, Two: "42.0"},
			{One: -128.0, Two: "-128.0"},
		}

		for _, tu := range tests {
			out := test.Eval(t, tu.Two)

			assert.Equal(t, object.FLOAT, out.Type())
			assert.Equal(t, tu.One, out.GoValue())
		}
	})

	t.Run("NullExpression", func(t *testing.T) {
		out := test.Eval(t, "nil")

		assert.Equal(t, object.NULL, out.Type())
		assert.Equal(t, nil, out.GoValue())
	})

	t.Run("BooleanExpression", func(t *testing.T) {
		tests := test.TestingTuple2[string, bool]{
			// EQ
			{One: "1 == 2", Two: false},
			{One: "1 == true", Two: true},
			{One: "true == 1", Two: true},
			// NEQ
			{One: "1 != 2", Two: true},
			// LT
			{One: "1 < 1", Two: false},
			{One: "1 < 2", Two: true},
			{One: "5 < 2", Two: false},
			// GT
			{One: "1 > 1", Two: false},
			{One: "1 > 2", Two: false},
			{One: "5 > 2", Two: true},
			// LT-EQ
			{One: "5 <= 5", Two: true},
			{One: "5 <= 4", Two: false},
			{One: "5 <= 6", Two: true},
			// GT-EQ
			{One: "5 >= 5", Two: true},
			{One: "5 >= 4", Two: true},
			{One: "5 >= 6", Two: false},
			// NOT
			{One: "!(false == 0)", Two: false},
		}

		tests.Each(func(statement string, value bool) {
			t.Run(statement, func(t *testing.T) {
				out := test.Eval(t, statement)

				assert.Equal(t, value, out.GoValue())
			})
		})

		t.Run("true", func(t *testing.T) {
			out := test.Eval(t, "true")

			assert.Equal(t, object.BOOLEAN, out.Type())
			assert.Equal(t, true, out.GoValue())
		})

		t.Run("false", func(t *testing.T) {
			out := test.Eval(t, "false")

			assert.Equal(t, object.BOOLEAN, out.Type())
			assert.Equal(t, false, out.GoValue())
		})
	})

	t.Run("Bang Prefix", func(t *testing.T) {
		tests := []test.Tuple2[bool, string]{
			{One: false, Two: "!true"},
			{One: true, Two: "!false"},
			{One: true, Two: "!!true"},
			{One: false, Two: "!!false"},
			{One: false, Two: "!1"},
		}

		for _, tu := range tests {
			out := test.Eval(t, tu.Two)

			assert.Equal(t, object.BOOLEAN, out.Type())
			assert.Equal(t, tu.One, out.GoValue())
		}
	})
}

func TestExecuteEdgeCases(t *testing.T) {
	t.Run("invalid stack frame id for function defined in while loop", func(t *testing.T) {
		source := `
		const indirect = fn(child) { child() }

		let i = 0;
		let indirect_called = false

		while (i < 2) {
			i = i + 1

			indirect(fn() { indirect_called = true })
		}

		return indirect_called;
		`

		result := test.Eval(t, source)

		assert.Equal(t, true, result.GoValue())
	})

	t.Run("change value in parent scope after function is defined", func(t *testing.T) {
		source := `
		let value = 1;

		const get_value = fn() { value }

		value = 2;

		return get_value()
		`

		result := test.Eval(t, source)

		assert.Equal(t, int64(2), result.GoValue())
	})
}
