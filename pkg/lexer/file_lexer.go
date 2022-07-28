package lexer

import (
	"io"
	"sort"
	"text/scanner"
	"unicode"

	"github.com/maddiesch/marble/pkg/token"
	"github.com/pkg/errors"
)

type FileLexer struct {
	s *scanner.Scanner

	keywords  map[string]token.Kind
	operators []Operator
}

func (f *FileLexer) RegisterOperator(op Operator) {
	f.operators = append(f.operators, op)

	sort.Slice(f.operators, func(i, j int) bool {
		return f.operators[i].pres > f.operators[j].pres
	})
}

func (f *FileLexer) RegisterSingleRuneOperator(r rune, k token.Kind) {
	f.RegisterOperator(
		NewSingleValueOperator(r, k, OperatorPrecedenceLow),
	)
}

func (f *FileLexer) RegisterDoubleRuneOperator(r1, r2 rune, k token.Kind) {
	f.RegisterOperator(
		NewDoubleValueOperator(r1, r2, k, OperatorPrecedenceHigh),
	)
}

func (f *FileLexer) RegisterKeyword(n string, k token.Kind) {
	f.keywords[n] = k
}

func (l *FileLexer) NextToken() token.Token {
	createToken := func(k token.Kind) token.Token {
		return token.Token{
			Kind:     k,
			Value:    l.s.TokenText(),
			Location: token.Location(l.s.Pos()),
		}
	}

	r := l.s.Scan()

	switch r {
	case scanner.EOF:
		return createToken(token.EndOfInput)
	case scanner.String:
		return createToken(token.String)
	case scanner.Int:
		return createToken(token.Integer)
	case scanner.Float:
		return createToken(token.Float)
	case scanner.Comment:
		return createToken(token.Comment)
	case scanner.Ident:
		if kw, ok := l.keywords[l.s.TokenText()]; ok {
			return createToken(kw)
		}
		return createToken(token.Identifier)
	default:
		for _, op := range l.operators {
			if op.r1 != r {
				continue
			}
			if op.r2 == operatorRuneUnused {
				return createToken(op.Kind)
			} else if op.r2 == l.s.Peek() {
				l.s.Scan() // Consume the next token

				return createToken(op.Kind)
			}
		}
		return createToken(token.Invalid)
	}
}

var _ Lexer = (*FileLexer)(nil)

func New(name string, in io.Reader) (Lexer, error) {
	if seeker, ok := in.(io.Seeker); ok {
		_, err := seeker.Seek(0, io.SeekStart)
		if err != nil {
			return nil, errors.Wrap(err, "failed to seek the reader to the start of the file")
		}
	}

	s := new(scanner.Scanner)
	s = s.Init(in)
	s.Filename = name
	s.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanStrings | scanner.ScanComments
	s.IsIdentRune = func(r rune, i int) bool {
		return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) && i > 0
	}

	lex := &FileLexer{
		s:         s,
		keywords:  make(map[string]token.Kind),
		operators: make([]Operator, 0, 16),
	}

	return lex.init()
}

func (f *FileLexer) init() (*FileLexer, error) {
	f.RegisterKeyword("let", token.Let)
	f.RegisterKeyword("const", token.Const)
	f.RegisterKeyword("fn", token.Function)
	f.RegisterKeyword("return", token.Return)
	f.RegisterKeyword("if", token.If)
	f.RegisterKeyword("else", token.Else)
	f.RegisterKeyword("not", token.Not)
	f.RegisterKeyword("nil", token.LiteralNull)
	f.RegisterKeyword("false", token.LiteralFalse)
	f.RegisterKeyword("true", token.LiteralTrue)
	f.RegisterKeyword("mut", token.Mutate)
	f.RegisterKeyword("do", token.Do)

	f.RegisterSingleRuneOperator('=', token.Assign)
	f.RegisterSingleRuneOperator('+', token.Plus)
	f.RegisterSingleRuneOperator('-', token.Minus)
	f.RegisterSingleRuneOperator('*', token.Asterisk)
	f.RegisterSingleRuneOperator('/', token.Slash)
	f.RegisterSingleRuneOperator('<', token.LessThan)
	f.RegisterSingleRuneOperator('>', token.GreaterThan)
	f.RegisterSingleRuneOperator('!', token.Bang)

	f.RegisterDoubleRuneOperator('=', '=', token.Equate)
	f.RegisterDoubleRuneOperator('!', '=', token.NotEquate)
	f.RegisterDoubleRuneOperator(':', ':', token.DoubleColon)

	f.RegisterSingleRuneOperator(',', token.Comma)
	f.RegisterSingleRuneOperator('?', token.Question)
	f.RegisterSingleRuneOperator(';', token.Semicolon)
	f.RegisterSingleRuneOperator(':', token.Colon)
	f.RegisterSingleRuneOperator('(', token.LParen)
	f.RegisterSingleRuneOperator(')', token.RParen)
	f.RegisterSingleRuneOperator('[', token.LBracket)
	f.RegisterSingleRuneOperator(']', token.RBracket)
	f.RegisterSingleRuneOperator('{', token.LBrace)
	f.RegisterSingleRuneOperator('}', token.RBrace)
	f.RegisterSingleRuneOperator('~', token.Tilde)
	f.RegisterSingleRuneOperator('^', token.Carrot)
	f.RegisterSingleRuneOperator('%', token.Percent)
	f.RegisterSingleRuneOperator('.', token.Period)

	return f, nil
}
