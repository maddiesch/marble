package lexer_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/maddiesch/marble/internal/test"
	"github.com/maddiesch/marble/pkg/core/lexer"
	"github.com/maddiesch/marble/pkg/core/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("given a standard input", func(t *testing.T) {
		source := strings.NewReader(`let message = "Hello World!";`)

		l, err := lexer.New(t.Name(), source)

		assert.NoError(t, err)
		assert.NotNil(t, l)
	})
}

func TestLexer(t *testing.T) {
	t.Run("FileLexer", func(t *testing.T) {
		t.Run("Tokens", func(t *testing.T) {
			tests := test.TestingTuple3[string, token.Kind, string]{
				{One: `1.0`, Two: token.Float, Three: "1.0"},
				{One: `42`, Two: token.Integer, Three: "42"},
				{One: "`", Two: token.Invalid, Three: "`"},
				{One: "// Comment", Two: token.Comment, Three: "// Comment"},
				{One: "while", Two: token.While, Three: "while"},
			}

			tests.Each(func(source string, expectedKind token.Kind, expectedValue string) {
				t.Run(source, func(t *testing.T) {
					l, err := lexer.New(t.Name(), strings.NewReader(source))

					if assert.NoError(t, err, "failed to create lexer") {
						first := l.NextToken()

						assert.Equal(t, expectedKind, first.Kind)
						assert.Equal(t, expectedValue, first.Value)
					}
				})
			})
		})

		t.Run("ArrayLiteral", func(t *testing.T) {
			tokens := ParseTokenSlice(t, strings.NewReader(`[1, "foo", 3.14]`))

			tokens.assert(t,
				token.LBracket,
				token.Integer,
				token.Comma,
				token.String,
				token.Comma,
				token.Float,
				token.RBracket,
			)
		})

		t.Run("keywords", func(t *testing.T) {
			source := strings.NewReader(`let const fn return if else not nil false true`)

			tokens := ParseTokenSlice(t, source)

			tokens.assert(t,
				token.Let,
				token.Const,
				token.Function,
				token.Return,
				token.If,
				token.Else,
				token.Not,
				token.LiteralNull,
				token.LiteralFalse,
				token.LiteralTrue,
			)
		})

		t.Run("operators", func(t *testing.T) {
			source := strings.NewReader(`= + - * / < > ! == != , ? ; : () [] {} ~ ^ %`)

			tokens := ParseTokenSlice(t, source)

			tokens.assert(t,
				token.Assign,
				token.Plus,
				token.Minus,
				token.Asterisk,
				token.Slash,
				token.LessThan,
				token.GreaterThan,
				token.Bang,
				token.Equate,
				token.NotEquate,
				token.Comma,
				token.Question,
				token.Semicolon,
				token.Colon,
				token.LParen,
				token.RParen,
				token.LBracket,
				token.RBracket,
				token.LBrace,
				token.RBrace,
				token.Tilde,
				token.Carrot,
				token.Percent,
			)
		})
	})

	t.Run("double rune token value", func(t *testing.T) {
		tests := []test.Tuple3[string, int, string]{
			{One: "foo == bar", Two: 1, Three: "=="},
			{One: "foo != bar", Two: 1, Three: "!="},
			{One: "foo >= bar", Two: 1, Three: ">="},
			{One: "let r = foo <= bar", Two: 4, Three: "<="},
		}

		for _, tt := range tests {
			tokens := ParseTokenSlice(t, strings.NewReader(tt.One))

			if assert.GreaterOrEqual(t, len(tokens), tt.Two+1, "not enough tokens") {
				assert.Equal(t, tt.Three, tokens[tt.Two].Value)
			}
		}
	})

	t.Run("given a basic assignment", func(t *testing.T) {
		source := strings.NewReader(`let message = "Hello World!";`)

		tokens := ParseTokenSlice(t, source)

		tokens.assert(t,
			token.Let,
			token.Identifier,
			token.Assign,
			token.String,
			token.Semicolon,
		)
	})
}

type TestI interface {
	assert.TestingT
	require.TestingT

	Name() string
}

type TokenSlice []token.Token

func (tt TokenSlice) assert(t assert.TestingT, kk ...token.Kind) bool {
	if len(kk) != len(tt) {
		return assert.Fail(t, fmt.Sprintf("TokenSlice and Expected tokens length do not match: %d vs %d", len(kk), len(tt)))
	}

	for i := 0; i < len(tt); i++ {
		if tt[i].Is(kk[i]) {
			continue
		}

		return assert.Fail(t, fmt.Sprintf("Unexpected token kind at index %d. Expected=%s, Actual=%s", i, kk[i], tt[i].Kind))
	}

	return true
}

func ParseTokenSlice(t TestI, source io.Reader) TokenSlice {
	l, err := lexer.New(t.Name(), source)

	require.NoError(t, err)

	tokens := lexer.ReadAll(l)

	return TokenSlice(tokens)
}
