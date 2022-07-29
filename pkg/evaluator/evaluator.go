package evaluator

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/ast"
	"github.com/maddiesch/marble/pkg/evaluator/runtime"
	"github.com/maddiesch/marble/pkg/object"
	"github.com/maddiesch/marble/pkg/object/math"
)

func Evaluate(node ast.Node) (object.Object, error) {
	switch node := node.(type) {
	case *ast.Program:
		return _evalStatementList(node.StatementList)
	case *ast.ExpressionStatement:
		return Evaluate(node.Expression)
	case *ast.IntegerExpression:
		return &object.Integer{Value: node.Value}, nil
	case *ast.FloatExpression:
		return &object.Floating{Value: node.Value}, nil
	case *ast.BooleanExpression:
		return &object.Boolean{Value: node.Value}, nil
	case *ast.NullExpression:
		return &object.Null{}, nil
	case *ast.PrefixExpression:
		return _evalPrefixExpression(node)
	case *ast.InfixExpression:
		return _evalInfixExpression(node)
	default:
		return &object.Void{}, nil
	}
}

func _evalInfixExpression(node *ast.InfixExpression) (object.Object, error) {
	lhs, err := Evaluate(node.Left)
	if err != nil {
		return nil, err
	}
	rhs, err := Evaluate(node.Right)
	if err != nil {
		return nil, err
	}

	return _evalInfixOperator(node, lhs, rhs)
}

func _evalInfixOperator(n *ast.InfixExpression, lhs, rhs object.Object) (object.Object, error) {
	if op, ok := math.OperatorFor(n.Operator); ok {
		if left, ok := lhs.(object.BasicArithmeticEvaluator); ok {
			return left.PerformBasicArithmeticOperation(op, rhs)
		}
	}
	return nil, runtime.NewError(
		runtime.TypeError,
		"Unable to perform infix operation for given types",
		runtime.ErrorValue("Operator", n.Operator),
		runtime.ErrorValue("Left", lhs),
		runtime.ErrorValue("Right", rhs),
		runtime.ErrorValue("Location", n.SourceToken().Location),
	)
}

func _evalPrefixExpression(node *ast.PrefixExpression) (object.Object, error) {
	right, err := Evaluate(node.Expression)
	if err != nil {
		return nil, err
	}
	switch node.Operator {
	case "!":
		return _evalBangForExpression(right)
	case "-":
		return _evalNegateExpression(right, node)
	default:
		return nil, fmt.Errorf("failed to evaluate prefix expression")
	}
}

func _evalNegateExpression(obj object.Object, n ast.Node) (object.Object, error) {
	switch o := obj.(type) {
	case *object.Integer:
		return &object.Integer{Value: -o.Value}, nil
	case *object.Floating:
		return &object.Floating{Value: -o.Value}, nil
	default:
		return nil, IllegalExpression{Node: n, Message: "Unable to negate given object"}
	}
}

func _evalBangForExpression(obj object.Object) (object.Object, error) {
	switch o := obj.(type) {
	case *object.Boolean:
		return &object.Boolean{Value: !o.Value}, nil
	default:
		return &object.Boolean{Value: false}, nil
	}
}

func _evalStatementList(list []ast.Statement) (object.Object, error) {
	var result object.Object
	var err error

	for _, node := range list {
		result, err = Evaluate(node)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
