package token_test

import (
	"testing"

	"github.com/maddiesch/marble/pkg/token"
	"github.com/stretchr/testify/assert"
)

func TestLocation(t *testing.T) {
	l := token.Location{
		Filename: "testing_file",
		Line:     42,
		Offset:   129,
		Column:   16,
	}

	t.Run("String()", func(t *testing.T) {
		assert.Equal(t, "testing_file:42+16", l.String())
	})
}

func TestToken(t *testing.T) {
	tok := &token.Token{
		Kind:  token.String,
		Value: `"Hello World!"`,
	}

	t.Run("Is()", func(t *testing.T) {
		assert.True(t, tok.Is(token.String))
		assert.False(t, tok.Is(token.Identifier))
	})

	t.Run("IsNot()", func(t *testing.T) {
		assert.False(t, tok.IsNot(token.String))
		assert.True(t, tok.IsNot(token.Identifier))
	})

	t.Run("String()", func(t *testing.T) {
		assert.Equal(t, `STRING("Hello World!")`, tok.String())
	})
}
