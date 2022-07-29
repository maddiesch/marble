package evaluator

import (
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
	case *ast.NegateExpression:
		return _evalNegateExpression(node)
	case *ast.NotExpression:
		return _evalNotExpression(node)
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

	switch n.Operator {
	case "<", "<=", ">", ">=":
		return _evalComparisonInfixExpression(n, lhs, rhs)
	case "==", "!=":
		return _evalBooleanResultInfixExpression(n, lhs, rhs)
	default:
		return nil, runtime.NewError(
			runtime.TypeError,
			"Unable to perform infix operation for given types",
			runtime.ErrorValue("Operator", n.Operator),
			runtime.ErrorValue("Left", lhs),
			runtime.ErrorValue("Right", rhs),
			runtime.ErrorValue("Location", n.SourceToken().Location),
		)
	}
}

func _evalComparisonInfixExpression(n *ast.InfixExpression, lhs, rhs object.Object) (object.Object, error) {
	comparable, ok := lhs.(object.ComparisionEvaluator)
	if !ok {
		return nil, runtime.NewError(
			runtime.TypeError,
			"Unable to perform comparison for the given type",
			runtime.ErrorValue("Operator", n.Operator),
			runtime.ErrorValue("Left", lhs),
			runtime.ErrorValue("Right", rhs),
			runtime.ErrorValue("Location", n.SourceToken().Location),
		)
	}

	lessThan, err := comparable.PerformLessThanComparison(rhs)
	if err != nil {
		return nil, err
	}

	equal, err := comparable.PerformEqualityCheck(rhs)
	if err != nil {
		return nil, err
	}

	switch n.Operator {
	case "<":
		return object.Bool(lessThan && !equal), nil
	case ">":
		return object.Bool(!lessThan && !equal), nil
	case "<=":
		return object.Bool(lessThan || equal), nil
	case ">=":
		return object.Bool(!lessThan || equal), nil
	default:
		panic("unable to handle the given operator, this is an interpreter error as the operator should not have been passed here!")
	}
}

func _evalBooleanResultInfixExpression(n *ast.InfixExpression, lhs, rhs object.Object) (object.Object, error) {
	equateable, ok := lhs.(object.EqualityEvaluator)
	if !ok {
		return nil, runtime.NewError(
			runtime.TypeError,
			"Unable to perform equality check for the given type",
			runtime.ErrorValue("FailureReason", "Expression left-hand side does not conform to EqualityEvaluator"),
			runtime.ErrorValue("Operator", n.Operator),
			runtime.ErrorValue("Left", lhs),
			runtime.ErrorValue("Right", rhs),
			runtime.ErrorValue("Location", n.SourceToken().Location),
		)
	}

	eq, err := equateable.PerformEqualityCheck(rhs)
	if err != nil {
		return nil, err
	}

	if n.Operator == "!=" {
		eq = !eq
	}

	return object.Bool(eq), nil
}

func _evalNotExpression(node *ast.NotExpression) (object.Object, error) {
	right, err := Evaluate(node.Expression)
	if err != nil {
		return nil, err
	}

	b := new(object.Boolean)

	if err := object.CoerceToType(right, b); err != nil {
		return nil, err
	}

	return object.Bool(!b.Value), nil
}

func _evalNegateExpression(n *ast.NegateExpression) (object.Object, error) {
	obj, err := Evaluate(n.Expression)
	if err != nil {
		return nil, err
	}
	if left, ok := obj.(object.BasicArithmeticEvaluator); ok {
		return left.PerformBasicArithmeticOperation(math.OperationMultiply, object.Int(-1))
	}
	return nil, runtime.NewError(
		runtime.TypeError,
		"Unable to negate expression",
		runtime.ErrorValue("Value", obj),
		runtime.ErrorValue("Node", n),
	)
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
