package object_test

import (
	"fmt"
	"testing"

	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/core/object"
	"github.com/stretchr/testify/assert"
)

func TestArithmeticVisitor(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		tests := test.TestingTuple3[object.Object, object.Object, any]{
			{One: object.NewInteger(42), Two: object.NewInteger(128), Three: int64(170)},
			{One: object.NewFloat(3.14), Two: object.NewFloat(3.14), Three: float64(6.28)},
			{One: object.NewString("Hello "), Two: object.NewString("World!"), Three: "Hello World!"},
		}

		tests.Each(func(left, right object.Object, expected any) {
			t.Run(fmt.Sprintf("%s + %s = %s", left.DebugString(), right.DebugString(), fmt.Sprint(expected)), func(t *testing.T) {
				result, err := object.GetArithmeticResult('+', left, right)

				if assert.Nil(t, err, "visitor contained an error") {
					assert.Equal(t, expected, result.GoValue())
				}
			})
		})
	})

	t.Run("Sub", func(t *testing.T) {
		tests := test.TestingTuple3[object.Object, object.Object, any]{
			{One: object.NewInteger(42), Two: object.NewInteger(128), Three: int64(-86)},
			{One: object.NewFloat(3.14), Two: object.NewFloat(3.14), Three: float64(0.0)},
		}

		tests.Each(func(left, right object.Object, expected any) {
			t.Run(fmt.Sprintf("%s - %s = %s", left.DebugString(), right.DebugString(), fmt.Sprint(expected)), func(t *testing.T) {
				result, err := object.GetArithmeticResult('-', left, right)

				if assert.Nil(t, err, "visitor contained an error") {
					assert.Equal(t, expected, result.GoValue())
				}
			})
		})
	})

	t.Run("Multiply", func(t *testing.T) {
		tests := test.TestingTuple3[object.Object, object.Object, any]{
			{One: object.NewInteger(42), Two: object.NewInteger(128), Three: int64(5376)},
			{One: object.NewFloat(3.14), Two: object.NewFloat(8.4), Three: float64(26.376)},
		}

		tests.Each(func(left, right object.Object, expected any) {
			t.Run(fmt.Sprintf("%s * %s = %s", left.DebugString(), right.DebugString(), fmt.Sprint(expected)), func(t *testing.T) {
				result, err := object.GetArithmeticResult('*', left, right)

				if assert.Nil(t, err, "visitor contained an error") {
					assert.Equal(t, expected, result.GoValue())
				}
			})
		})
	})

	t.Run("Divide", func(t *testing.T) {
		tests := test.TestingTuple3[object.Object, object.Object, any]{
			{One: object.NewInteger(100), Two: object.NewInteger(20), Three: int64(5)},
			{One: object.NewFloat(18.9), Two: object.NewFloat(3.0), Three: float64(6.3)},
		}

		tests.Each(func(left, right object.Object, expected any) {
			t.Run(fmt.Sprintf("%s / %s = %s", left.DebugString(), right.DebugString(), fmt.Sprint(expected)), func(t *testing.T) {
				result, err := object.GetArithmeticResult('/', left, right)

				if assert.Nil(t, err, "visitor contained an error") {
					assert.Equal(t, expected, result.GoValue())
				}
			})
		})
	})
}
