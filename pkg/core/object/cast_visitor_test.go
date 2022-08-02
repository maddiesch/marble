package object_test

import (
	"fmt"
	"testing"

	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/core/object"
	"github.com/stretchr/testify/assert"
)

func TestCastVisitor(t *testing.T) {
	tests := test.TestingTuple3[object.Object, object.ObjectType, object.Object]{
		{One: object.NewInteger(42), Two: object.INTEGER, Three: object.NewInteger(42)},
		{One: object.NewInteger(42), Two: object.FLOAT, Three: object.NewFloat(42.0)},
		{One: object.NewInteger(1), Two: object.BOOLEAN, Three: object.NewBool(true)},
		{One: object.NewInteger(0), Two: object.BOOLEAN, Three: object.NewBool(false)},
		{One: object.NewBool(true), Two: object.INTEGER, Three: object.NewInteger(1)},
		{One: object.NewBool(false), Two: object.INTEGER, Three: object.NewInteger(0)},
		{One: object.NewBool(true), Two: object.FLOAT, Three: object.NewFloat(1.0)},
		{One: object.NewBool(false), Two: object.FLOAT, Three: object.NewFloat(0.0)},
		{One: object.NewBool(true), Two: object.BOOLEAN, Three: object.NewBool(true)},
		{One: object.NewBool(false), Two: object.BOOLEAN, Three: object.NewBool(false)},
		{One: object.NewFloat(42.0), Two: object.INTEGER, Three: object.NewInteger(42)},
		{One: object.NewFloat(42.0), Two: object.FLOAT, Three: object.NewFloat(42.0)},
		{One: object.NewReturn(object.NewInteger(343)), Two: object.FLOAT, Three: object.NewFloat(343.0)},
		{One: object.NewNull(), Two: object.NULL, Three: object.NewNull()},
		{One: object.NewNull(), Two: object.BOOLEAN, Three: object.NewBool(false)},
	}

	tests.Each(func(given object.Object, target object.ObjectType, outcome object.Object) {
		t.Run(fmt.Sprintf("Cast %s to %s", given.Type(), target), func(t *testing.T) {
			v := &object.CastVisitor{Target: target}

			given.Accept(v)

			result, err := v.Take()

			if assert.Nil(t, err, "visitor contained an error") {
				if assert.NotNil(t, result, "visitor did not assign a result value") {
					assert.Equal(t, outcome.Type(), result.Type())
					assert.Equal(t, outcome.DebugString(), result.DebugString())
				}
			}
		})
	})
}
