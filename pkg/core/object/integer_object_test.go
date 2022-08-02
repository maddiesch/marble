package object_test

import (
	"fmt"
	"testing"

	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/core/object"
	"github.com/stretchr/testify/assert"
)

func TestInteger(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		tests := []test.Tuple3[*object.IntegerObject, object.Object, bool]{
			{One: object.NewInteger(1), Two: object.NewInteger(1), Three: true},
			{One: object.NewInteger(1), Two: object.NewInteger(3), Three: false},
			{One: object.NewInteger(1), Two: object.NewFloat(1.0), Three: true},
			{One: object.NewInteger(1), Two: object.NewBool(true), Three: true},
		}

		for _, tu := range tests {
			t.Run(fmt.Sprintf("%s == %s = %t", tu.One.DebugString(), tu.Two.DebugString(), tu.Three), func(t *testing.T) {
				eq, err := object.GetObjectEquality(tu.One, tu.Two)

				if assert.Nil(t, err, "failed to perform equality check") {
					assert.Equal(t, tu.Three, eq)
				}
			})
		}
	})
}
