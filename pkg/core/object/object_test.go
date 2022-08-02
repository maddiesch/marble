package object_test

import (
	"fmt"
	"testing"

	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/core/object"
	"github.com/stretchr/testify/assert"
)

func TestObjectType(t *testing.T) {
	tests := test.TestingTuple2[object.Object, object.ObjectType]{
		{One: &object.Boolean{}, Two: object.BOOLEAN},
		{One: &object.Floating{}, Two: object.FLOAT},
		{One: &object.Integer{}, Two: object.INTEGER},
		{One: &object.Null{}, Two: object.NULL},
		{One: &object.Void{}, Two: object.VOID},
		{One: &object.StringObject{}, Two: object.STRING},
	}

	tests.Each(func(o object.Object, ot object.ObjectType) {
		t.Run(fmt.Sprintf("validate object type %s", o), func(t *testing.T) {
			assert.Equal(t, ot, o.Type())
		})
	})
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
