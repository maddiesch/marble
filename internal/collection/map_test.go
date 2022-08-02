package collection_test

import (
	"testing"

	"github.com/maddiesch/marble/internal/collection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapSlice(t *testing.T) {
	s := []string{"foo", "bar", "baz"}
	o := collection.MapSlice(s, func(v string) string {
		return v + ".mapped"
	})

	require.Equal(t, 3, len(o))

	assert.Equal(t, "foo.mapped", o[0])
	assert.Equal(t, "bar.mapped", o[1])
	assert.Equal(t, "baz.mapped", o[2])
}
