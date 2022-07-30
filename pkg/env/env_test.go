package env_test

import (
	"testing"

	"github.com/maddiesch/marble/pkg/env"
	"github.com/maddiesch/marble/pkg/object"
	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	e := env.New()

	e.Set("foo", object.Int(1), true)
	e.Set("bar", object.String("baz"), true)

	e.Push()

	e.Set("baz", object.Bool(true), true)
	e.Set("foo", object.Bool(false), true)

	baz, ok := e.Get("baz")
	if assert.True(t, ok) {
		assert.Equal(t, true, baz.GoValue())
	}

	bar, ok := e.Get("bar")
	if assert.True(t, ok) {
		assert.Equal(t, "baz", bar.GoValue())
	}

	e.Pop()

	_, ok = e.Get("baz")

	assert.False(t, ok)

	foo, ok := e.Get("foo")
	if assert.True(t, ok) {
		assert.Equal(t, int64(1), foo.GoValue())
	}
}
