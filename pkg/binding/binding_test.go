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

		assert.False(t, b1.GetState("bar", true).IsSet(), "Binding 1 should not have 'bar' set")
		assert.True(t, b2.GetState("foo", true).IsSet(), "Binding 2 should have 'foo' set")
		assert.False(t, b2.GetState("foo", false).IsSet(), "Binding 2 should have 'foo' set without recursive searching")
		assert.False(t, b2.GetState("foo", true).IsCurrent(), "Binding 2 should not have 'foo' as the current frame")
		assert.False(t, b2.GetState("foo", true).IsMutable(), "Binding 2 should not have 'foo' as a mutable value")
		assert.True(t, b2.GetState("bar", true).IsMutable(), "Binding 2 should allow 'bar' to be mutated")
	})
}
