// Package lexer handles parsing the source code passed as a reader and
// converting it to a list of tokens.
package lexer

import (
	"github.com/maddiesch/marble/pkg/core/token"
)

type Lexer interface {
	NextToken() token.Token
}

type Operator struct {
	token.Kind

	pres OperatorPrecedence

	r1 rune
	r2 rune
}

func NewSingleValueOperator(r rune, k token.Kind, p OperatorPrecedence) Operator {
	return Operator{
		Kind: k,
		pres: p,
		r1:   r,
		r2:   operatorRuneUnused,
	}
}

func NewDoubleValueOperator(r1, r2 rune, k token.Kind, p OperatorPrecedence) Operator {
	return Operator{
		Kind: k,
		pres: p,
		r1:   r1,
		r2:   r2,
	}
}

type OperatorPrecedence uint8

const (
	OperatorPrecedenceLow  OperatorPrecedence = 0
	OperatorPrecedenceHigh OperatorPrecedence = 255
)

const (
	operatorRuneUnused rune = -1
)

// ReadAll reads all Lexer tokens into a slice and returns that slice.
//
// It will continue to read until a EndOfInput token in found. Once that token
// is reached, it will return all collected tokens. The EndOfInput token will
// not be included in the returned slice.
func ReadAll(l Lexer) []token.Token {
	result := make([]token.Token, 0, 64)

ReadLoop:
	for {
		t := l.NextToken()

		if t.Is(token.EndOfInput) {
			break ReadLoop
		}

		result = append(result, t)
	}

	return result
}
