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
		{One: object.NewBool(false), Two: object.BOOLEAN},
		{One: object.NewFloat(1.0), Two: object.FLOAT},
		{One: object.NewInteger(0), Two: object.INTEGER},
		{One: object.NewNull(), Two: object.NULL},
		{One: object.NewVoid(), Two: object.VOID},
		{One: &object.StringObject{}, Two: object.STRING},
	}

	tests.Each(func(o object.Object, ot object.ObjectType) {
		t.Run(fmt.Sprintf("validate object type %s", o), func(t *testing.T) {
			assert.Equal(t, ot, o.Type())
		})
	})
}
