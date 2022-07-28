package parser

import (
	"github.com/maddiesch/marble/pkg/ast"
	"github.com/maddiesch/marble/pkg/token"
)

// Run executes the parser
func (p *Parser) Run() *ast.Program {
	defer untrace(trace("Run"))

	p.parseErr = make([]error, 0)

	program := &ast.Program{
		StatementList: make([]ast.Statement, 0),
	}

	for p.currentToken.IsNot(token.EndOfInput) {
		statement, err := p.parseCurrentTokenIntoStatement()
		if err != nil {
			p.encounteredErr(err)
		}
		if statement != nil {
			program.StatementList = append(program.StatementList, statement)
		}
		p.advance()
	}

	return program
}

func (p *Parser) parseCurrentTokenIntoStatement() (ast.Statement, error) {
	defer untrace(trace("parseCurrentTokenIntoStatement"))

	switch p.currentToken.Kind {
	case token.Invalid:
		return nil, InvalidTokenError{Token: p.currentToken}
	case token.Let:
		return p.parseLetStatement()
	case token.Const:
		return p.parseConstantStatement()
	case token.Mutate:
		return p.parseMutateStatement()
	case token.Return:
		return p.parseReturnStatement()
	case token.Comment:
		return nil, nil
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() (ast.Statement, error) {
	defer untrace(trace("parseReturnStatement"))

	stmt := &ast.ReturnStatement{
		Token: p.currentToken,
	}
	p.advance()

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

func (p *Parser) _parseAssignmentStatement() (*ast.IdentifierExpression, ast.Expression, error) {
	defer untrace(trace("_parseAssignmentStatement"))

	if !p.currentTokenIs(token.Identifier) {
		return nil, nil, UnexpectedTokenError{Token: p.nextToken, Expected: token.Identifier}
	}

	identifier := &ast.IdentifierExpression{
		Token: p.currentToken,
		Value: p.currentToken.Value,
	}

	p.advance() // Consume the identifier

	if !p.currentTokenIs(token.Assign) {
		return nil, nil, UnexpectedTokenError{Token: p.nextToken, Expected: token.Assign}
	}
	p.advance() // consume the assignment

	expr, err := p.parseExpression(Lowest)
	if err != nil {
		return nil, nil, err
	}

	if p.nextTokenIs(token.Semicolon) {
		p.advance()
	}

	return identifier, expr, nil
}

func (p *Parser) parseLetStatement() (ast.Statement, error) {
	defer untrace(trace("parseLetStatement"))

	current := p.currentToken
	p.advance()
	identifier, expression, err := p._parseAssignmentStatement()
	if err != nil {
		return nil, err
	}

	return &ast.LetStatement{
		Token:      current,
		Identifier: identifier,
		Expression: expression,
	}, nil
}

func (p *Parser) parseConstantStatement() (ast.Statement, error) {
	defer untrace(trace("parseConstantStatement"))

	current := p.currentToken
	p.advance()
	identifier, expression, err := p._parseAssignmentStatement()
	if err != nil {
		return nil, err
	}

	return &ast.ConstantStatement{
		Token:      current,
		Identifier: identifier,
		Expression: expression,
	}, nil
}

func (p *Parser) parseMutateStatement() (ast.Statement, error) {
	defer untrace(trace("parseMutateStatement"))

	current := p.currentToken
	p.advance()
	identifier, expression, err := p._parseAssignmentStatement()
	if err != nil {
		return nil, err
	}

	return &ast.MutateStatement{
		Token:      current,
		Identifier: identifier,
		Expression: expression,
	}, nil
}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	defer untrace(trace("parseBlockStatement"))

	block := &ast.BlockStatement{
		Token:         p.currentToken,
		StatementList: make([]ast.Statement, 0),
	}

	p.advance() // Consume LBrace

	for !p.currentTokenIs(token.RBrace) {
		if p.currentTokenIs(token.EndOfInput) {
			return nil, UnexpectedTokenError{
				Token:    p.currentToken,
				Expected: token.RBrace,
			}
		}
		stmt, err := p.parseCurrentTokenIntoStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			block.StatementList = append(block.StatementList, stmt)
		}
		p.advance()
	}

	return block, nil
}
