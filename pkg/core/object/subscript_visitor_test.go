package object_test

import (
	"fmt"
	"testing"

	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/core/object"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptVisitor(t *testing.T) {
	tests := test.TestingTuple2[object.Object, object.Object]{
		{One: object.NewInteger(0), Two: object.NewInteger(18)},
	}

	testArray := object.NewArray([]object.Object{
		object.NewInteger(18),
	})

	tests.Each(func(key object.Object, expected object.Object) {
		t.Run(fmt.Sprintf("Subscript %s to %s", key.Type(), fmt.Sprint(key.GoValue())), func(t *testing.T) {
			result, err := object.GetSubscriptValue(testArray, key)

			if assert.Nil(t, err, "visitor contained an error") {
				assert.Equal(t, expected.DebugString(), result.DebugString())
			}
		})
	})
}
