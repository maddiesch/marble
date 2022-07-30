package stack_test

import (
	"testing"

	"github.com/maddiesch/marble/pkg/stack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
	t.Run("At", func(t *testing.T) {
		s := stack.New("foo", "bar", "baz")

		assert.Equal(t, "foo", s.At(0))
		assert.Equal(t, "bar", s.At(1))
		assert.Equal(t, "baz", s.At(2))

		assert.Panics(t, func() {
			s.At(5)
		})
	})

	t.Run("Search", func(t *testing.T) {
		s := stack.New("foo", "bar", "baz")

		t.Run("when found", func(t *testing.T) {
			i, v := s.Search(func(v string) bool {
				return v == "bar"
			})

			assert.Equal(t, 1, i)
			assert.Equal(t, "bar", v)
		})

		t.Run("when not found", func(t *testing.T) {
			i, _ := s.Search(func(v string) bool {
				return v == "none"
			})

			assert.Equal(t, -1, i)
		})
	})

	t.Run("RevIter", func(t *testing.T) {
		s := stack.New("foo", "bar")

		for i, v := range s.RevIter() {
			switch i {
			case 0:
				assert.Equal(t, "bar", v)
			case 1:
				assert.Equal(t, "foo", v)
			default:
				assert.Fail(t, "unexpected iterator index")
			}
		}
	})

	t.Run("Iter", func(t *testing.T) {
		s := stack.New("foo", "bar")

		for i, v := range s.Iter() {
			switch i {
			case 0:
				assert.Equal(t, "foo", v)
			case 1:
				assert.Equal(t, "bar", v)
			default:
				assert.Fail(t, "unexpected iterator index")
			}
		}
	})

	t.Run("Last", func(t *testing.T) {
		s := stack.New("foo", "bar")

		f, _ := s.Last()

		assert.Equal(t, "bar", f)

		t.Run("when empty", func(t *testing.T) {
			s = stack.New[string]()

			_, ok := s.Last()

			assert.False(t, ok)
		})
	})

	t.Run("First", func(t *testing.T) {
		s := stack.New("foo", "bar")

		f, _ := s.First()

		assert.Equal(t, "foo", f)

		t.Run("when empty", func(t *testing.T) {
			s = stack.New[string]()

			_, ok := s.First()

			assert.False(t, ok)
		})
	})

	t.Run("Pop", func(t *testing.T) {
		s := stack.New("foo", "bar")

		assert.Equal(t, "bar", s.Pop())

		assert.Equal(t, 1, s.Len())

		t.Run("when empty", func(t *testing.T) {
			s = stack.New[string]()

			assert.Panics(t, func() {
				s.Pop()
			})
		})
	})

	t.Run("Push", func(t *testing.T) {
		s := stack.New[string]()

		s.Push("foo")
		s.Push("bar")

		require.Equal(t, 2, s.Len())
	})

	t.Run("Len", func(t *testing.T) {
		s := stack.New("foo", "bar")

		assert.Equal(t, 2, s.Len())
	})
}
