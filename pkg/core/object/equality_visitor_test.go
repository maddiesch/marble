package object_test

import (
	"fmt"
	"testing"

	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/core/object"
	"github.com/stretchr/testify/assert"
)

func TestEqualityVisitor(t *testing.T) {
	tests := test.TestingTuple3[object.Object, object.Object, bool]{
		{One: object.NewInteger(42), Two: object.NewInteger(42), Three: true},
		{One: object.NewInteger(42), Two: object.NewInteger(0), Three: false},
		{One: object.NewInteger(42), Two: object.NewFloat(42.0), Three: true},
		{One: object.NewFloat(42.0), Two: object.NewFloat(42.0), Three: true},
		{One: object.NewFloat(42.0), Two: object.NewFloat(100.0), Three: false},
		{One: object.NewFloat(42.0), Two: object.NewInteger(42), Three: true},
		{One: object.NewBool(true), Two: object.NewBool(true), Three: true},
		{One: object.NewBool(false), Two: object.NewBool(true), Three: false},
		{One: object.NewBool(true), Two: object.NewInteger(1), Three: true},
	}

	tests.Each(func(left object.Object, right object.Object, expected bool) {
		testName := fmt.Sprintf("Equal (%s) %s == (%s) %s = %t", left.Type(), fmt.Sprint(left.GoValue()), right.Type(), fmt.Sprint(right.GoValue()), expected)
		t.Run(testName, func(t *testing.T) {
			equal, err := object.GetObjectEquality(left, right)

			if assert.Nil(t, err, "visitor contained an error") {
				assert.Equal(t, expected, equal)
			}
		})
	})
}
