package evaluator_test

import (
	"testing"

	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/core/evaluator"
	"github.com/maddiesch/marble/pkg/core/object"
	"github.com/stretchr/testify/require"
)

var _benchResult object.Object

func BenchmarkFib10(b *testing.B) {
	b.StopTimer()

	source := `const fib = fn(n) { if (n < 2) { return n } else { return fib(n-1) + fib(n-2) } }; fib(10);`
	program := test.CreateProgram(b, source)

	// Run the program once outside of the benchmark loop to validate the program source
	binding := evaluator.NewBinding()
	result, err := evaluator.Evaluate(binding, program)

	require.NoError(b, err)
	require.Equal(b, int64(55), result.GoValue())

	b.StartTimer()

	var r object.Object

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			binding := evaluator.NewBinding()
			r, _ = evaluator.Evaluate(binding, program)
		}
	})

	// Ensure the compiler doesn't optimize the loop away
	_benchResult = r
}
