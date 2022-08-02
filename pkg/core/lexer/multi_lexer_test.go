package lexer_test

import (
	"strings"
	"testing"

	"github.com/maddiesch/marble/pkg/core/lexer"
	"github.com/maddiesch/marble/pkg/core/token"
	"github.com/stretchr/testify/assert"
)

func TestMultiLexer(t *testing.T) {
	s1 := strings.NewReader("foo")
	s2 := strings.NewReader("bar")

	l1, _ := lexer.New("multi-1", s1)
	l2, _ := lexer.New("multi-2", s2)

	multi := lexer.NewLexer(l1, l2)

	tokens := lexer.ReadAll(multi)

	if assert.Equal(t, 2, len(tokens)) {
		assert.Equal(t, token.Identifier, tokens[0].Kind)
		assert.Equal(t, "foo", tokens[0].Value)
		assert.Equal(t, "multi-1", tokens[0].Location.Filename)

		assert.Equal(t, token.Identifier, tokens[1].Kind)
		assert.Equal(t, "bar", tokens[1].Value)
		assert.Equal(t, "multi-2", tokens[1].Location.Filename)
	}
}
