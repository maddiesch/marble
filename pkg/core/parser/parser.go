// Package parser converts lexers tokens into an AST
package parser

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/maddiesch/marble/pkg/core/ast"
	"github.com/maddiesch/marble/pkg/core/lexer"
	"github.com/maddiesch/marble/pkg/core/token"
)

type (
	ParserPrefixFunc func() (ast.Expression, error)
	ParserInfixFunc  func(ast.Expression) (ast.Expression, error)
)

type Parser struct {
	lexer lexer.Lexer

	currentToken token.Token
	nextToken    token.Token

	parseErr []error

	precedence map[token.Kind]Precedence

	prefixParser map[token.Kind]ParserPrefixFunc
	infixParser  map[token.Kind]ParserInfixFunc
}

func (p *Parser) Err() error {
	if len(p.parseErr) == 0 {
		return nil
	}
	return &ParseError{Children: p.parseErr}
}

func New(l lexer.Lexer) *Parser {
	p := &Parser{
		lexer: l,
	}

	return p.init()
}

func (p *Parser) init() *Parser {
	p.precedence = map[token.Kind]Precedence{
		token.Equate:           Equals,
		token.NotEquate:        Equals,
		token.LessThan:         LessGreater,
		token.LessThanEqual:    LessGreater,
		token.GreaterThan:      LessGreater,
		token.GreaterThanEqual: LessGreater,
		token.Plus:             Sum,
		token.Minus:            Sum,
		token.Slash:            Product,
		token.Asterisk:         Product,
		token.LParen:           Call,
		token.DoubleColon:      Call,
		token.Period:           Call,
		token.LBracket:         Subscript,
		token.Assign:           Assignment,
	}

	p.prefixParser = make(map[token.Kind]ParserPrefixFunc)
	p.infixParser = make(map[token.Kind]ParserInfixFunc)

	p.registerPrefixParser(token.Identifier, p.parseIdentifierExpression)
	p.registerPrefixParser(token.Integer, p.parseIntegerExpression)
	p.registerPrefixParser(token.Float, p.parseFloatExpression)
	p.registerPrefixParser(token.String, p.parseStringExpression)
	p.registerPrefixParser(token.Bang, p.parsePrefixExpression)
	p.registerPrefixParser(token.Not, p.parsePrefixExpression)
	p.registerPrefixParser(token.Minus, p.parsePrefixExpression)
	p.registerPrefixParser(token.LiteralNull, p.parseNullExpression)
	p.registerPrefixParser(token.LiteralTrue, p.parseBooleanExpression)
	p.registerPrefixParser(token.LiteralFalse, p.parseBooleanExpression)
	p.registerPrefixParser(token.LParen, p.parseGroupedExpression)
	p.registerPrefixParser(token.If, p.parseIfExpression)
	p.registerPrefixParser(token.While, p.parseWhileExpression)
	p.registerPrefixParser(token.Function, p.parseFunctionExpression)
	p.registerPrefixParser(token.Do, p.parseDoExpression)
	p.registerPrefixParser(token.Defined, p.parseDefinedExpression)
	p.registerPrefixParser(token.LBracket, p.parseBracketExpression)

	p.registerInfixParser(token.Plus, p.parseInfixExpression)
	p.registerInfixParser(token.Minus, p.parseInfixExpression)
	p.registerInfixParser(token.Slash, p.parseInfixExpression)
	p.registerInfixParser(token.Asterisk, p.parseInfixExpression)
	p.registerInfixParser(token.Equate, p.parseInfixExpression)
	p.registerInfixParser(token.NotEquate, p.parseInfixExpression)
	p.registerInfixParser(token.LessThan, p.parseInfixExpression)
	p.registerInfixParser(token.LessThanEqual, p.parseInfixExpression)
	p.registerInfixParser(token.GreaterThan, p.parseInfixExpression)
	p.registerInfixParser(token.GreaterThanEqual, p.parseInfixExpression)
	p.registerInfixParser(token.Assign, p.parseAssignInfixExpression)
	p.registerInfixParser(token.LParen, p.parseCallExpression)
	p.registerInfixParser(token.DoubleColon, p.parseDoubleColonExpression)
	p.registerInfixParser(token.Period, p.parsePeriodExpression)
	p.registerInfixParser(token.LBracket, p.parseSubscriptExpression)

	// Advance 2 tokens to pre-load the next & current tokens
	p.advance()
	p.advance()

	return p
}

func (p *Parser) advance() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

func (p *Parser) currentTokenIs(k token.Kind) bool {
	return p.currentToken.Is(k)
}

func (p *Parser) nextTokenIs(k token.Kind) bool {
	return p.nextToken.Is(k)
}

func (p *Parser) encounteredErr(e error) {
	if os.Getenv("MARBLE_ABORT_FAST") == "true" {
		spew.Dump(p)
		panic(e)
	}
	p.parseErr = append(p.parseErr, e)
}

func (p *Parser) registerPrefixParser(t token.Kind, fn ParserPrefixFunc) {
	p.prefixParser[t] = fn
}

func (p *Parser) registerInfixParser(t token.Kind, fn ParserInfixFunc) {
	p.infixParser[t] = fn
}

func (p *Parser) precedenceFor(t token.Token) Precedence {
	if pre, ok := p.precedence[t.Kind]; ok {
		return pre
	}
	return Lowest
}
