package object_test

import (
	"fmt"
	"testing"

	"github.com/maddiesch/marble/pkg/object"
	"github.com/stretchr/testify/assert"
)

func TestObjectType(t *testing.T) {
	tests := []struct {
		object   object.Object
		expected object.ObjectType
	}{
		{&object.Boolean{}, object.BOOLEAN},
		{&object.Floating{}, object.FLOAT},
		{&object.Integer{}, object.INTEGER},
		{&object.Null{}, object.NULL},
		{&object.Void{}, object.VOID},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("validate object type %s", test.expected), func(t *testing.T) {
			assert.Equal(t, test.expected, test.object.Type())
		})
	}
}

func TestChainCoerceTo(t *testing.T) {
	t.Run("float -> int -> bool", func(t *testing.T) {
		start := object.Float(1.0)

		b, err := object.ChainCoerceTo(start, object.INTEGER, object.BOOLEAN)
		if assert.NoError(t, err) {
			assert.Equal(t, true, b.GoValue())
		}
	})
}
