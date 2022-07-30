package evaluator

import (
	"fmt"
	"runtime/debug"

	"github.com/maddiesch/marble/internal/collection"
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
	case *ast.StatementExpression:
		return Evaluate(env, node.Statement)
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
		return _evalAssignmentNode(env, node, false)
	case *ast.MutateStatement:
		return _evalAssignmentNode(env, node, true)
	case *ast.ConstantStatement:
		return _evalAssignmentNode(env, node, false)
	case *ast.DeleteStatement:
		return _evalDeleteStatement(env, node)
	case *ast.ArrayExpression:
		return _evalArrayNode(env, node)
	case *ast.SubscriptExpression:
		return _evalSubscriptNode(env, node)
	case *ast.DefinedExpression:
		return object.Bool(env.StateFor(node.Identifier.Value, false).Defined()), nil
	case *ast.WhileExpression:
		return _evalWhileExpression(env, node)
	case *ast.BreakStatement:
		return object.Instruction(object.NativeInstructionBreak), nil
	case *ast.ContinueStatement:
		return object.Instruction(object.NativeInstructionContinue), nil
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
	case *ast.FunctionExpression:
		return _evalFunctionExpression(env, node)
	case *ast.CallExpression:
		return _evalCallExpression(env, node)
	default:
		return nil, runtime.NewError(
			runtime.InterpreterError,
			"I haven't been taught how to interpret this node yet!",
			runtime.ErrorValue("NodeType", fmt.Sprintf("%T", node)),
			runtime.ErrorValue("Node", node),
		)
	}
}

func _evalWhileExpression(e *env.Env, node *ast.WhileExpression) (object.Object, error) {
	var result object.Object
	result = &object.Null{}

EvalLoop:
	for {
		condition, err := Evaluate(e, node.Condition)
		if err != nil {
			return nil, err
		}
		boolean, err := object.CoerceTo(condition, object.BOOLEAN)
		if err != nil {
			return nil, err
		}

		if !boolean.(*object.Boolean).Value {
			break
		}

		blockResult, err := _evalStatementList(e, node.Block.StatementList, true, false)
		if err != nil {
			return nil, err
		}

		switch r := blockResult.(type) {
		case *object.ReturnObject:
			return r.Value, nil
		case *object.NativeInstruction:
			switch r.IType {
			case object.NativeInstructionBreak:
				break EvalLoop
			case object.NativeInstructionContinue:
				continue EvalLoop
			}
		}

		result = blockResult
	}

	return result, nil
}

func _evalSubscriptNode(e *env.Env, node *ast.SubscriptExpression) (object.Object, error) {
	rec, err := Evaluate(e, node.Receiver)
	if err != nil {
		return nil, err
	}

	sub, ok := rec.(object.SubscriptEvaluator)
	if !ok {
		return nil, runtime.NewError(runtime.ArgumentError,
			"Subscript receiver does not implement subscript accessor",
			runtime.ErrorValue("Receiver", rec),
		)
	}

	val, err := Evaluate(e, node.Value)
	if err != nil {
		return nil, err
	}

	return sub.Subscript(val)
}

func _evalArrayNode(e *env.Env, node *ast.ArrayExpression) (object.Object, error) {
	array := make([]object.Object, len(node.Elements))

	for i, n := range node.Elements {
		if v, err := Evaluate(e, n); err != nil {
			return nil, err
		} else {
			array[i] = v
		}
	}

	return object.Array(array), nil
}

func _evalNativeFunction(e *env.Env, fn *object.NativeFunctionObject, node *ast.CallExpression) (object.Object, error) {
	if len(node.Arguments) != fn.ArgumentCount {
		return nil, runtime.NewError(runtime.ArgumentError,
			"Unexpected number of arguments in function call",
			runtime.ErrorValue("Expected", fn.ArgumentCount),
			runtime.ErrorValue("Received", len(node.Arguments)),
		)
	}

	arguments := make([]object.Object, len(node.Arguments))
	for i, arg := range node.Arguments {
		val, err := Evaluate(e, arg)
		if err != nil {
			return nil, err
		}
		arguments[i] = val
	}

	e.PushTo(1) // Always push to start frame
	defer e.Restore()

	e.Push()
	defer e.Pop()

	return fn.Body(e, arguments)
}

func _evalCallExpression(e *env.Env, node *ast.CallExpression) (object.Object, error) {
	fn, err := Evaluate(e, node.Function)
	if err != nil {
		return nil, err
	}

	switch closure := fn.(type) {
	case *object.ClosureLiteral:
		return _evalCallClosureLiteral(e, closure, node)
	case *object.NativeFunctionObject:
		return _evalNativeFunction(e, closure, node)
	default:
		return nil, runtime.NewError(runtime.CallError,
			"Unable to call non-function type",
			runtime.ErrorValue("Function", fn),
		)
	}
}

func _evalCallClosureLiteral(e *env.Env, closure *object.ClosureLiteral, node *ast.CallExpression) (object.Object, error) {
	if len(node.Arguments) != len(closure.ParameterList) {
		return nil, runtime.NewError(runtime.ArgumentError,
			"Unexpected number of arguments in function call",
			runtime.ErrorValue("Expected", len(closure.ParameterList)),
			runtime.ErrorValue("Received", len(node.Arguments)),
		)
	}

	// We need to evaluate the arguments before we push to the closure's frame
	// We also can't bind the arguments into the environment until we have pushed into the closure's frame
	// So we loop once to evaluate the arguments, push to the closure's frame, then bind the argument values into the frame
	arguments := make([]object.Object, len(node.Arguments))
	for i, arg := range node.Arguments {
		val, err := Evaluate(e, arg)
		if err != nil {
			return nil, err
		}
		arguments[i] = val
	}

	if !e.PushTo(closure.FrameID) {
		fmt.Println(e.DebugString())
		fmt.Println(node.String())
		debug.PrintStack()
		return nil, runtime.NewError(runtime.InterpreterError,
			"Unable to push to the expected execution state for function call!",
			runtime.ErrorValue("StackFrame", closure.FrameID),
		)
	}
	defer e.Restore()

	e.Push()
	defer e.Pop()

	for i, val := range arguments {
		name := closure.ParameterList[i]
		e.Set(name, val, false)
	}

	return _evalStatementList(e, closure.Body.StatementList, true, true)
}

func _evalFunctionExpression(e *env.Env, node *ast.FunctionExpression) (*object.ClosureLiteral, error) {
	parameters := collection.MapSlice(node.Parameters, func(l *ast.IdentifierExpression) string {
		return l.Value
	})
	return object.Closure(parameters, node.BlockStatement, e.CurrentFrame()), nil
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

func _evalAssignmentNode(e *env.Env, node ast.AssignmentStatement, update bool) (object.Object, error) {
	name := node.Label()

	currentState := e.StateFor(name, node.CurrentFrameOnly())

	if node.RequireDefined() && !currentState.Defined() {
		return nil, runtime.NewError(runtime.UnknownIdentifierError,
			"Unable to assing to an undefined variable",
			runtime.ErrorValue("Name", name),
			runtime.ErrorValue("Location", node.SourceToken().Location),
		)
	}

	if currentState.Defined() {
		if !currentState.Mutable() {
			return nil, runtime.NewError(runtime.ConstantError,
				"Attempting to assign a new value to a constant",
				runtime.ErrorValue("Name", name),
				runtime.ErrorValue("Location", node.SourceToken().Location),
			)
		}

		if node.RequireUndefined() {
			return nil, runtime.NewError(runtime.InterpreterError,
				"Unable to perform assignment on a defined variable",
				runtime.ErrorValue("Name", name),
				runtime.ErrorValue("Location", node.SourceToken().Location),
			)
		}
	}

	value, err := Evaluate(e, node.AssignmentExpression())

	if err != nil {
		return nil, err
	}
	if update {
		e.Update(name, value, node.CurrentFrameOnly())
	} else {
		e.Set(name, value, node.Mutable())
	}
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

	// Handle Concatenation
	if lConcat, ok := lhs.(object.ConcatingEvaluator); ok && n.Operator == "+" {
		return lConcat.Concat(rhs)
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

	result = &object.Null{}

	for _, node := range list {
		result, err = Evaluate(env, node)
		if err != nil {
			return nil, err
		}

		switch r := result.(type) {
		case *object.ReturnObject:
			if unwrapReturn {
				return r.Value, nil
			} else {
				return r, nil
			}
		case *object.NativeInstruction:
			switch r.IType {
			case object.NativeInstructionBreak, object.NativeInstructionContinue:
				return r, nil
			}
		}
	}

	return result, nil
}
