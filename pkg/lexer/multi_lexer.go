package lexer

import "github.com/maddiesch/marble/pkg/token"

type MultiLexer struct {
	children     []Lexer
	currentChild int
}

var _ Lexer = (*MultiLexer)(nil)

func (l *MultiLexer) NextToken() token.Token {
ChildToken:
	if len(l.children) <= l.currentChild {
		return token.Token{
			Kind: token.EndOfInput,
		}
	}
	c := l.children[l.currentChild]

	t := c.NextToken()

	if t.Is(token.EndOfInput) {
		l.currentChild += 1
		goto ChildToken
	}

	return t
}

func NewLexer(c ...Lexer) *MultiLexer {
	return &MultiLexer{children: c}
}

func (l *MultiLexer) Add(c Lexer) {
	l.children = append(l.children, c)
}
