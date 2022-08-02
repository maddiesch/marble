package object_test

import (
	"testing"

	"github.com/maddiesch/marble/pkg/core/object"
	"github.com/stretchr/testify/assert"
)

func TestCopyObject(t *testing.T) {
	start := object.NewInteger(4839)
	dest := object.NewInteger(0)

	t.Run("given a pointer source", func(t *testing.T) {
		object.CopyObject(start, dest)

		assert.Equal(t, int64(4839), dest.GoValue())
	})
}
