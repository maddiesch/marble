package parser_test

import (
	"testing"

	"github.com/maddiesch/marble/pkg/parser"
	"github.com/maddiesch/marble/pkg/token"
	"github.com/stretchr/testify/assert"
)

func TestInvalidTokenError(t *testing.T) {
	token := token.Token{
		Kind: token.Invalid,
		Location: token.Location{
			Filename: "testing",
			Offset:   1,
			Line:     2,
			Column:   3,
		},
	}

	t.Run("Error string format", func(t *testing.T) {
		e := &parser.InvalidTokenError{
			Token: token,
		}

		expected := "InvalidToken: Parser encountered an invalid token\n  Source: testing:2+3"

		assert.Equal(t, expected, e.Error())
	})
}
