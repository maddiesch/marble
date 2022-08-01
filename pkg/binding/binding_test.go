package binding_test

import (
	"testing"

	"github.com/maddiesch/marble/pkg/binding"
	"github.com/stretchr/testify/assert"
)

func TestBinding(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		b1 := binding.New[string](nil)
		b1.Set("foo", "bar", binding.F_CONST)

		b2 := binding.New(b1)

		b2.Set("bar", "baz", 0)
		b2.Set("baz", "foo", 0)

		assert.False(t, b1.GetState("bar", true).IsSet())
		assert.True(t, b2.GetState("foo", true).IsSet())
		assert.False(t, b2.GetState("foo", true).IsCurrent())
		assert.False(t, b2.GetState("foo", true).IsMutable())
		assert.True(t, b2.GetState("bar", true).IsMutable())
	})
}
