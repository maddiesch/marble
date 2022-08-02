package object_test

import (
	"fmt"
	"testing"

	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/core/object"
	"github.com/stretchr/testify/assert"
)

func TestComparisonVisitor(t *testing.T) {
	tests := test.TestingTuple3[object.Object, object.Object, bool]{
		{One: object.NewInteger(42), Two: object.NewInteger(50), Three: true},
		{One: object.NewInteger(42), Two: object.NewInteger(42), Three: false},
		{One: object.NewInteger(42), Two: object.NewFloat(22.0), Three: false},
		{One: object.NewString("a"), Two: object.NewString("z"), Three: true},
	}

	tests.Each(func(left object.Object, right object.Object, expected bool) {
		testName := fmt.Sprintf("LessThan (%s) %s < (%s) %s = %t", left.Type(), fmt.Sprint(left.GoValue()), right.Type(), fmt.Sprint(right.GoValue()), expected)
		t.Run(testName, func(t *testing.T) {
			lessThan, err := object.GetLessThanComparison(left, right)

			if assert.Nil(t, err, "visitor contained an error") {
				assert.Equal(t, expected, lessThan)
			}
		})
	})
}
