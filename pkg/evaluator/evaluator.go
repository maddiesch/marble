package evaluator

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/ast"
	"github.com/maddiesch/marble/pkg/object"
	"golang.org/x/exp/constraints"
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
	switch n.Operator {
	case "+", "-", "*", "/":
		return _evalBasicArithmetic(n, lhs, rhs)
	default:
		return nil, IllegalExpression{Node: n, Message: "Illegal infix operator"}
	}
}

func _evalBasicArithmetic(n *ast.InfixExpression, lhs, rhs object.Object) (object.Object, error) {
	cast, ok := rhs.Cast(lhs.Type())
	if !ok {
		return nil, ObjectCastError{
			Node:   n,
			Object: rhs,
			Target: lhs.Type(),
		}
	}

	switch lhs := lhs.(type) {
	case *object.Integer:
		rhs := cast.(*object.Integer)
		return &object.Integer{Value: _execMathOp(n.Operator, lhs.Value, rhs.Value)}, nil
	case *object.Floating:
		rhs := cast.(*object.Floating)
		return &object.Floating{Value: _execMathOp(n.Operator, lhs.Value, rhs.Value)}, nil
	default:
		return nil, IllegalExpression{Node: n, Message: "Unable to perform arithmetic operation"}
	}
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

func _execMathOp[V constraints.Integer | constraints.Float](o string, l, r V) V {
	switch o {
	case "+":
		return l + r
	case "-":
		return l - r
	case "/":
		return l / r
	case "*":
		return l * r
	default:
		panic("illegal math operator " + o)
	}
}

func _execValueComparable[T constraints.Ordered](o string, l, r T) bool {
	switch o {
	case "<":
		return l < r
	case "<=":
		return l <= r
	case ">":
		return l > r
	case ">=":
		return l >= r
	default:
		panic("illegal comparision operator " + o)
	}
}
