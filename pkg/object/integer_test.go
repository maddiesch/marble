package object_test

import (
	"fmt"
	"testing"

	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/object"
	"github.com/maddiesch/marble/pkg/object/math"
	"github.com/stretchr/testify/assert"
)

func TestInteger(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		tests := []test.Tuple3[*object.Integer, object.Object, bool]{
			{One: object.Int(1), Two: object.Int(1), Three: true},
			{One: object.Int(1), Two: object.Int(3), Three: false},
			{One: object.Int(1), Two: object.Float(1.0), Three: true},
			{One: object.Int(1), Two: object.Bool(true), Three: true},
		}

		for _, tu := range tests {
			t.Run(fmt.Sprintf("%s == %s = %t", tu.One.Description(), tu.Two.Description(), tu.Three), func(t *testing.T) {
				eq, ok := tu.One.Equal(tu.Two)

				if assert.True(t, ok, "failed to perform equality check") {
					assert.Equal(t, tu.Three, eq)
				}
			})
		}
	})

	t.Run("LessThan", func(t *testing.T) {
		tests := []test.Tuple3[*object.Integer, object.Object, bool]{
			{One: object.Int(1), Two: object.Int(2), Three: true},
			{One: object.Int(1), Two: object.Int(0), Three: false},
			{One: object.Int(1), Two: object.Int(1), Three: false},
			{One: object.Int(1), Two: object.Float(2.0), Three: true},
			{One: object.Int(0), Two: object.Bool(true), Three: true},
		}

		for _, tu := range tests {
			t.Run(fmt.Sprintf("%s < %s = %t", tu.One.Description(), tu.Two.Description(), tu.Three), func(t *testing.T) {
				eq, ok := tu.One.LessThan(tu.Two)

				if assert.True(t, ok, "failed to perform equality check") {
					assert.Equal(t, tu.Three, eq)
				}
			})
		}
	})

	t.Run("Add", func(t *testing.T) {
		tests := []test.Tuple3[*object.Integer, math.ArithmeticAddition, *object.Integer]{
			{One: object.Int(1), Two: object.Int(2), Three: object.Int(3)},
			{One: object.Int(1), Two: object.Float(2.0), Three: object.Int(3)},
			{One: object.Int(6), Two: object.Int(-8), Three: object.Int(-2)},
		}

		for _, tu := range tests {
			t.Run(fmt.Sprintf("%s + %s = %s", tu.One.Description(), tu.Two.(object.Object).Description(), tu.Three.Description()), func(t *testing.T) {
				eq, ok := tu.One.Add(tu.Two)

				if assert.True(t, ok, "failed to perform equality check") {
					assert.Equal(t, tu.Three, eq)
				}
			})
		}
	})

	t.Run("Subtract", func(t *testing.T) {
		tests := []test.Tuple3[*object.Integer, math.ArithmeticSubtraction, *object.Integer]{
			{One: object.Int(1), Two: object.Int(2), Three: object.Int(-1)},
			{One: object.Int(1), Two: object.Float(2.0), Three: object.Int(-1)},
			{One: object.Int(6), Two: object.Int(-8), Three: object.Int(14)},
		}

		for _, tu := range tests {
			t.Run(fmt.Sprintf("%s - %s = %s", tu.One.Description(), tu.Two.(object.Object).Description(), tu.Three.Description()), func(t *testing.T) {
				eq, ok := tu.One.Sub(tu.Two)

				if assert.True(t, ok, "failed to perform equality check") {
					assert.Equal(t, tu.Three, eq)
				}
			})
		}
	})

	t.Run("Multiply", func(t *testing.T) {
		tests := []test.Tuple3[*object.Integer, math.ArithmeticMultiplication, *object.Integer]{
			{One: object.Int(1), Two: object.Int(2), Three: object.Int(2)},
			{One: object.Int(3), Two: object.Float(2.0), Three: object.Int(6)},
			{One: object.Int(6), Two: object.Int(-8), Three: object.Int(-48)},
		}

		for _, tu := range tests {
			t.Run(fmt.Sprintf("%s * %s = %s", tu.One.Description(), tu.Two.(object.Object).Description(), tu.Three.Description()), func(t *testing.T) {
				eq, ok := tu.One.Multiply(tu.Two)

				if assert.True(t, ok, "failed to perform equality check") {
					assert.Equal(t, tu.Three, eq)
				}
			})
		}
	})

	t.Run("Divide", func(t *testing.T) {
		tests := []test.Tuple3[*object.Integer, math.ArithmeticDivision, *object.Integer]{
			{One: object.Int(50), Two: object.Int(5), Three: object.Int(10)},
			{One: object.Int(48), Two: object.Float(8.0), Three: object.Int(6)},
			{One: object.Int(256), Two: object.Int(-8), Three: object.Int(-32)},
		}

		for _, tu := range tests {
			t.Run(fmt.Sprintf("%s * %s = %s", tu.One.Description(), tu.Two.(object.Object).Description(), tu.Three.Description()), func(t *testing.T) {
				eq, ok := tu.One.Divide(tu.Two)

				if assert.True(t, ok, "failed to perform equality check") {
					assert.Equal(t, tu.Three, eq)
				}
			})
		}
	})
}
