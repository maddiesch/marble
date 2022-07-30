package evaluator

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/ast"
	"github.com/maddiesch/marble/pkg/env"
	"github.com/maddiesch/marble/pkg/evaluator/runtime"
	"github.com/maddiesch/marble/pkg/object"
	"github.com/maddiesch/marble/pkg/object/math"
)

func Evaluate(env *env.Env, node ast.Node) (object.Object, error) {
	env.PushEval(node)
	defer env.PopEval()

	switch node := node.(type) {
	case *ast.Program:
		return _evalStatementList(env, node.StatementList, false, true)
	case *ast.BlockStatement:
		return _evalStatementList(env, node.StatementList, true, false)
	case *ast.ExpressionStatement:
		return Evaluate(env, node.Expression)
	case *ast.IntegerExpression:
		return object.Int(node.Value), nil
	case *ast.FloatExpression:
		return object.Float(node.Value), nil
	case *ast.BooleanExpression:
		return object.Bool(node.Value), nil
	case *ast.StringExpression:
		return object.String(node.Value), nil
	case *ast.NullExpression:
		return &object.Null{}, nil
	case *ast.NegateExpression:
		return _evalNegateExpression(env, node)
	case *ast.NotExpression:
		return _evalNotExpression(env, node)
	case *ast.InfixExpression:
		return _evalInfixExpression(env, node)
	case *ast.IfExpression:
		return _evalIfExpression(env, node)
	case *ast.ReturnStatement:
		return _evalReturnStatement(env, node)
	case *ast.LetStatement:
		return _evalAssignmentNode(env, node)
	case *ast.ConstantStatement:
		return _evalAssignmentNode(env, node)
	case *ast.DeleteStatement:
		return _evalDeleteStatement(env, node)
	case *ast.DefinedExpression:
		return object.Bool(env.StateFor(node.Identifier.Value, false).Defined()), nil
	case *ast.IdentifierExpression:
		value, ok := env.Get(node.Value)
		if !ok {
			return nil, runtime.NewError(
				runtime.UnknownIdentifierError,
				"Undefined identifier",
				runtime.ErrorValue("Identifier", node.Value),
				runtime.ErrorValue("Location", node.Token.Location),
			)
		}
		return value, nil
	default:
		return nil, runtime.NewError(
			runtime.InterpreterError,
			"I haven't been taught how to interpret this node yet!",
			runtime.ErrorValue("NodeType", fmt.Sprintf("%T", node)),
			runtime.ErrorValue("Node", node),
		)
	}
}

func _evalDeleteStatement(e *env.Env, node *ast.DeleteStatement) (object.Object, error) {
	name := node.Identifier.Value
	state := e.StateFor(name, false)
	switch state {
	case env.LabelStateUnassigned:
		return nil, runtime.NewError(runtime.UnknownIdentifierError,
			"Can't delete an undefined identifier",
			runtime.ErrorValue("Location", node.SourceToken().Location),
			runtime.ErrorValue("Label", name),
		)
	case env.LabelStateAssignedProtected, env.LabelStateAssignedImmutable:
		return nil, runtime.NewError(runtime.ConstantError,
			"Can't delete a constant identifier",
			runtime.ErrorValue("Location", node.SourceToken().Location),
			runtime.ErrorValue("Label", name),
		)
	case env.LabelStateAssignedMutable:
		val, ok := e.Get(name)
		if !ok {
			return nil, runtime.NewError(runtime.InterpreterError,
				"Unable to find an existing value for the label",
				runtime.ErrorValue("Location", node.SourceToken().Location),
				runtime.ErrorValue("Label", name),
			)
		}

		e.Delete(name, false)

		return val, nil
	default:
		panic("there is a state for the environment entry that we weren't expecting. This should not be possible!")
	}
}

func _evalAssignmentNode(env *env.Env, node ast.AssignmentStatement) (object.Object, error) {
	name := node.Label()
	if !env.StateFor(name, true).Mutable() {
		return nil, runtime.NewError(
			runtime.ConstantError,
			"Attempting to assign a new value to a constant",
			runtime.ErrorValue("Location", node.SourceToken().Location),
			runtime.ErrorValue("Label", name),
		)
	}
	value, err := Evaluate(env, node.AssignmentExpression())
	if err != nil {
		return nil, err
	}
	env.Set(name, value, node.Mutable())
	return value, nil
}

func _evalReturnStatement(env *env.Env, node *ast.ReturnStatement) (object.Object, error) {
	val, err := Evaluate(env, node.Expression)
	if err != nil {
		return nil, err
	}
	return object.Return(val), nil
}

func _evalIfExpression(env *env.Env, node *ast.IfExpression) (object.Object, error) {
	condition, err := Evaluate(env, node.Condition)
	if err != nil {
		return nil, err
	}

	boolean := object.Bool(false)
	if err := object.CoerceToType(condition, boolean); err != nil {
		return nil, err
	}

	if boolean.Value {
		return Evaluate(env, node.TrueStatement)
	} else if node.FalseStatement != nil {
		return Evaluate(env, node.FalseStatement)
	} else {
		return &object.Null{}, nil
	}
}

func _evalInfixExpression(env *env.Env, node *ast.InfixExpression) (object.Object, error) {
	lhs, err := Evaluate(env, node.Left)
	if err != nil {
		return nil, err
	}
	rhs, err := Evaluate(env, node.Right)
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

func _evalNotExpression(env *env.Env, node *ast.NotExpression) (object.Object, error) {
	right, err := Evaluate(env, node.Expression)
	if err != nil {
		return nil, err
	}

	b := new(object.Boolean)

	if err := object.CoerceToType(right, b); err != nil {
		return nil, err
	}

	return object.Bool(!b.Value), nil
}

func _evalNegateExpression(env *env.Env, n *ast.NegateExpression) (object.Object, error) {
	obj, err := Evaluate(env, n.Expression)
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

// We need to know if we should return the return object or unwrap to the real
// value. This is so that nested return statements will bubble up to the parent
// to handle the return statement
func _evalStatementList(env *env.Env, list []ast.Statement, pushStack, unwrapReturn bool) (object.Object, error) {
	if pushStack {
		env.Push()
		defer env.Pop()
	}
	var result object.Object
	var err error

	for _, node := range list {
		result, err = Evaluate(env, node)
		if err != nil {
			return nil, err
		}

		if ret, ok := result.(*object.ReturnObject); ok {
			if unwrapReturn {
				return ret.Value, nil
			} else {
				return ret, nil
			}
		}
	}

	return result, nil
}
