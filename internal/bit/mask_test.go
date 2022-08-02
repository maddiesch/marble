package bit_test

import (
	"testing"

	"github.com/maddiesch/marble/internal/bit"
	"github.com/stretchr/testify/assert"
)

func TestHas(t *testing.T) {
	b := uint8(1) << 2
	b |= 1 << 4

	assert.False(t, bit.Has(b, 1<<1))
	assert.True(t, bit.Has(b, 1<<2))
	assert.False(t, bit.Has(b, 1<<3))
	assert.True(t, bit.Has(b, 1<<4))
}
