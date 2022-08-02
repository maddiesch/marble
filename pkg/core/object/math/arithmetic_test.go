package math_test

import (
	"testing"

	"github.com/maddiesch/marble/pkg/core/object/math"
	"github.com/stretchr/testify/assert"
)

func TestOperatorFor(t *testing.T) {
	op, ok := math.OperatorFor("+")
	assert.True(t, ok)
	assert.Equal(t, math.OperationAdd, op)
}
