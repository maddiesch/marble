package parser

import (
	"strconv"

	"github.com/maddiesch/marble/pkg/ast"
	"github.com/maddiesch/marble/pkg/token"
)

func (p *Parser) parseExpressionStatement() (ast.Statement, error) {
	defer untrace(trace("parseExpressionStatement"))

	stmt := &ast.ExpressionStatement{
		Token: p.currentToken,
	}

	expr, err := p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	stmt.Expression = expr

	if p.nextTokenIs(token.Semicolon) {
		p.advance()
	}

	return stmt, nil
}

func (p *Parser) parseExpression(pre Precedence) (ast.Expression, error) {
	defer untrace(trace("parseExpression"))

	prefix := p.prefixParser[p.currentToken.Kind]
	if prefix == nil {
		return nil, ExpressionParseError{
			Token:   p.currentToken,
			Message: "No prefix parser found for expression",
		}
	}
	leftExp, err := prefix()
	if err != nil {
		return nil, err
	}

	for !p.nextTokenIs(token.Semicolon) && pre < p.precedenceFor(p.nextToken) {
		infix := p.infixParser[p.nextToken.Kind]
		if infix == nil {
			return leftExp, nil
		}

		p.advance()

		leftExp, err = infix(leftExp)
		if err != nil {
			return nil, err
		}
	}

	return leftExp, nil
}

func (p *Parser) parseIdentifierExpression() (ast.Expression, error) {
	defer untrace(trace("parseIdentifierExpression"))

	return &ast.IdentifierExpression{
		Token: p.currentToken,
		Value: p.currentToken.Value,
	}, nil
}

func (p *Parser) parseFloatExpression() (ast.Expression, error) {
	defer untrace(trace("parseFloatExpression"))

	value, err := strconv.ParseFloat(p.currentToken.Value, 64)
	if err != nil {
		return nil, LiteralParseError{
			Message:         "Failed to parse float literal",
			UnderlyingError: err,
			Token:           p.currentToken,
		}
	}
	return &ast.FloatExpression{
		Token: p.currentToken,
		Value: value,
	}, nil
}

func (p *Parser) parseIntegerExpression() (ast.Expression, error) {
	defer untrace(trace("parseIntegerExpression"))

	value, err := strconv.ParseInt(p.currentToken.Value, 0, 64)
	if err != nil {
		return nil, LiteralParseError{
			Message:         "Failed to parse integer literal",
			UnderlyingError: err,
			Token:           p.currentToken,
		}
	}
	return &ast.IntegerExpression{
		Token: p.currentToken,
		Value: value,
	}, nil
}

func (p *Parser) parseStringExpression() (ast.Expression, error) {
	defer untrace(trace("parseStringExpression"))

	value, err := strconv.Unquote(p.currentToken.Value)
	if err != nil {
		return nil, LiteralParseError{
			Message:         "Failed to parse string literal",
			UnderlyingError: err,
			Token:           p.currentToken,
		}
	}
	return &ast.StringExpression{
		Token: p.currentToken,
		Value: value,
	}, nil
}

func (p *Parser) parseBooleanExpression() (ast.Expression, error) {
	defer untrace(trace("parseBooleanExpression"))

	value, err := strconv.ParseBool(p.currentToken.Value)
	if err != nil {
		return nil, LiteralParseError{
			Message:         "Failed to parse boolean literal",
			UnderlyingError: err,
			Token:           p.currentToken,
		}
	}
	return &ast.BooleanExpression{
		Token: p.currentToken,
		Value: value,
	}, nil
}

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	defer untrace(trace("parsePrefixExpression"))

	op := p.currentToken.Value
	if op == "not" {
		op = "!"
	}
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: op,
	}

	p.advance()

	child, err := p.parseExpression(Prefix)
	if err != nil {
		return nil, err
	}

	expression.Expression = child

	return expression, nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	defer untrace(trace("parseInfixExpression"))

	expr := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Value,
		Left:     left,
	}

	pre := p.precedenceFor(p.currentToken)
	p.advance()

	right, err := p.parseExpression(pre)
	if err != nil {
		return nil, err
	}

	expr.Right = right

	return expr, nil
}

func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	defer untrace(trace("parseGroupedExpression"))

	p.advance()

	exp, err := p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	p.advance()

	if !p.currentTokenIs(token.RParen) {
		return nil, UnexpectedTokenError{
			Token:    p.currentToken,
			Expected: token.RParen,
		}
	}

	return exp, nil
}

func (p *Parser) parseIfExpression() (ast.Expression, error) {
	defer untrace(trace("parseIfExpression"))

	expr := &ast.IfExpression{
		Token: p.currentToken,
	}

	if !p.nextTokenIs(token.LParen) {
		return nil, UnexpectedTokenError{
			Token:    p.nextToken,
			Expected: token.LParen,
		}
	}
	p.advance() // Consume IF

	condition, err := p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	if !p.currentTokenIs(token.RParen) {
		return nil, UnexpectedTokenError{
			Token:    p.currentToken,
			Expected: token.RParen,
		}
	}
	p.advance() // Consume )

	expr.Condition = condition

	if !p.currentTokenIs(token.LBrace) {
		return nil, UnexpectedTokenError{
			Token:    p.currentToken,
			Expected: token.LBrace,
		}
	}

	trueBlock, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	expr.TrueStatement = trueBlock

	if p.nextTokenIs(token.Else) {
		p.advance() // Consume }
		p.advance() // Consume ELSE

		if !p.currentTokenIs(token.LBrace) {
			return nil, UnexpectedTokenError{
				Token:    p.currentToken,
				Expected: token.LBrace,
			}
		}

		falseBlock, err := p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
		p.advance() // Consume }

		expr.FalseStatement = falseBlock
	}

	return expr, nil
}

func (p *Parser) parseFunctionExpression() (ast.Expression, error) {
	defer untrace(trace("parseFunctionExpression"))

	fn := &ast.FunctionExpression{
		Token: p.currentToken,
	}

	if !p.nextTokenIs(token.LParen) {
		return nil, UnexpectedTokenError{
			Token:    p.nextToken,
			Expected: token.LParen,
		}
	}
	p.advance() // Consume FN

	params, err := p.parseParameterList()
	if err != nil {
		return nil, err
	}

	p.advance() // Consume RParen

	fn.Parameters = params

	if !p.currentTokenIs(token.LBrace) {
		return nil, UnexpectedTokenError{
			Token:    p.currentToken,
			Expected: token.LBrace,
		}
	}

	block, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	fn.BlockStatement = block

	return fn, nil
}

func (p *Parser) parseParameterList() ([]*ast.IdentifierExpression, error) {
	defer untrace(trace("parseParameterList"))

	result := make([]*ast.IdentifierExpression, 0)

	if p.nextTokenIs(token.RParen) {
		p.advance()
		return result, nil
	}

	p.advance()

	assignCurrentIdentifier := func() error {
		if !p.currentTokenIs(token.Identifier) {
			return UnexpectedTokenError{
				Token:    p.currentToken,
				Expected: token.Identifier,
			}
		}

		result = append(result, &ast.IdentifierExpression{
			Token: p.currentToken,
			Value: p.currentToken.Value,
		})

		return nil
	}

	if err := assignCurrentIdentifier(); err != nil {
		return nil, err
	}

	for p.nextTokenIs(token.Comma) {
		p.advance() // Current is now comma
		p.advance() // Advance to the next ident

		if err := assignCurrentIdentifier(); err != nil {
			return nil, err
		}
	}

	if !p.nextTokenIs(token.RParen) {
		return nil, UnexpectedTokenError{
			Token:    p.nextToken,
			Expected: token.RParen,
		}
	}

	p.advance() // Consume the RParen

	return result, nil
}

func (p *Parser) parseCallExpression(fn ast.Expression) (ast.Expression, error) {
	defer untrace(trace("parseCallExpression"))

	call := &ast.CallExpression{
		Token:    p.currentToken,
		Function: fn,
	}

	arguments, err := p.parseCallArguments()
	if err != nil {
		return nil, err
	}

	call.Arguments = arguments

	return call, nil
}

func (p *Parser) parseCallArguments() ([]ast.Expression, error) {
	defer untrace(trace("parseCallArguments"))

	result := make([]ast.Expression, 0)

	if p.nextTokenIs(token.RParen) {
		p.advance()
		return result, nil
	}

	p.advance()

	assignCurrentExpression := func() error {
		expr, err := p.parseExpression(Lowest)
		if err != nil {
			return err
		}

		result = append(result, expr)

		return nil
	}

	if err := assignCurrentExpression(); err != nil {
		return nil, err
	}

	for p.nextTokenIs(token.Comma) {
		p.advance() // Current is now comma
		p.advance() // Advance to the next expression

		if err := assignCurrentExpression(); err != nil {
			return nil, err
		}
	}

	if !p.nextTokenIs(token.RParen) {
		return nil, UnexpectedTokenError{
			Token:    p.nextToken,
			Expected: token.RParen,
		}
	}

	p.advance() // Consume the RParen

	return result, nil
}

func (p *Parser) parseDoubleColonExpression(left ast.Expression) (ast.Expression, error) {
	defer untrace(trace("parseDoubleColonExpression"))

	id, ok := left.(*ast.IdentifierExpression)

	if !ok {
		return nil, ExpressionParseError{
			Token:   p.currentToken,
			Message: "Double colon is required to be preceded by an identifier.",
		}
	}
	mod := &ast.ModuleExpression{
		Token:  p.currentToken,
		Module: id,
	}

	p.advance() // Consume Double Colon

	expr, err := p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	mod.Expression = expr

	return mod, nil
}

func (p *Parser) parsePeriodExpression(left ast.Expression) (ast.Expression, error) {
	defer untrace(trace("parsePeriodExpression"))

	chain := &ast.ChainExpression{
		Token: p.currentToken,
		From:  left,
	}

	p.advance() // Consume Period
	expr, err := p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}

	chain.To = expr

	return chain, nil
}

func (p *Parser) parseDoExpression() (ast.Expression, error) {
	doToken := p.currentToken

	p.advance() // Consume Do
	if !p.currentTokenIs(token.LBrace) {
		return nil, UnexpectedTokenError{
			Token:    p.currentToken,
			Expected: token.LBrace,
		}
	}

	block, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &ast.DoExpression{
		Token: doToken,
		Block: block,
	}, nil
}
