package object_test

import (
	"fmt"
	"testing"

	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/object"
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
				eq, err := tu.One.PerformEqualityCheck(tu.Two)

				if assert.NoError(t, err, "failed to perform equality check") {
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
				eq, err := tu.One.PerformLessThanComparison(tu.Two)

				if assert.NoError(t, err, "failed to perform equality check") {
					assert.Equal(t, tu.Three, eq)
				}
			})
		}
	})
}
