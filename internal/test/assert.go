package test

import (
	"github.com/maddiesch/marble/pkg/core/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertLiteralValue(t TestingT, expect any, expr ast.Expression) bool {
	switch val := expr.(type) {
	case *ast.IdentifierExpression:
		return assert.Equal(t, expect, val.Value)
	case *ast.IntegerExpression:
		return assert.Equal(t, expect, val.Value)
	default:
		assert.Fail(t, "Unexpected Literal Expression", "Found %T", expr)
		return false
	}
}

func AssertInfixExpression(t TestingT, exp ast.Expression, left any, op string, right any) bool {
	infix, ok := exp.(*ast.InfixExpression)
	if !ok {
		assert.Fail(t, "Invalid Type", "Expected an InfixExpression, got %T", exp)
		return false
	}

	if !assert.Equal(t, op, infix.Operator) {
		return false
	}

	if !AssertLiteralValue(t, left, infix.Left) {
		return false
	}

	return AssertLiteralValue(t, right, infix.Right)
}

func AssertIdentifier(t TestingT, name string, e ast.Expression) bool {
	if id, ok := e.(*ast.IdentifierExpression); ok {
		return assert.Equal(t, name, id.Value)
	}
	assert.Fail(t, "Expected identifier expression", "expected *ast.IdentifierExpression, got %T", e)
	return false
}

func AssertStatementType(t TestingT, exp ast.Statement, val ast.Statement) bool {
	return assert.IsType(t, exp, val, "Unexpected statement type Expected=%T Got=%T", exp, val)
}

func RequireType[T any](t TestingT, exp T, val any) T {
	require.IsType(t, exp, val)
	return val.(T)
}
