package evaluator

import (
	"fmt"

	"github.com/maddiesch/marble/internal/collection"
	"github.com/maddiesch/marble/pkg/ast"
	"github.com/maddiesch/marble/pkg/binding"
	"github.com/maddiesch/marble/pkg/evaluator/runtime"
	"github.com/maddiesch/marble/pkg/native"
	"github.com/maddiesch/marble/pkg/object"
	"github.com/maddiesch/marble/pkg/object/math"
)

func NewBinding() *binding.Binding[object.Object] {
	b := binding.New[object.Object](nil)
	native.Bind(b)
	return b
}

func Evaluate(b *binding.Binding[object.Object], node ast.Node) (object.Object, error) {
	switch node := node.(type) {
	case *ast.Program:
		return _evalStatementList(b, node.StatementList, false, true)
	case *ast.BlockStatement:
		return _evalStatementList(b, node.StatementList, true, false)
	case *ast.ExpressionStatement:
		return Evaluate(b, node.Expression)
	case *ast.DoExpression:
		return _evalStatementList(b, node.Block.StatementList, true, true)
	case *ast.StatementExpression:
		return Evaluate(b, node.Statement)
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
		return _evalNegateExpression(b, node)
	case *ast.NotExpression:
		return _evalNotExpression(b, node)
	case *ast.InfixExpression:
		return _evalInfixExpression(b, node)
	case *ast.IfExpression:
		return _evalIfExpression(b, node)
	case *ast.ReturnStatement:
		return _evalReturnStatement(b, node)
	case *ast.LetStatement:
		return _evalAssignmentNode(b, node, false)
	case *ast.MutateStatement:
		return _evalAssignmentNode(b, node, true)
	case *ast.ConstantStatement:
		return _evalAssignmentNode(b, node, false)
	case *ast.DeleteStatement:
		return _evalDeleteStatement(b, node)
	case *ast.ArrayExpression:
		return _evalArrayNode(b, node)
	case *ast.SubscriptExpression:
		return _evalSubscriptNode(b, node)
	case *ast.DefinedExpression:
		return object.Bool(b.GetState(node.Identifier.Value, false).IsSet()), nil
	case *ast.WhileExpression:
		return _evalWhileExpression(b, node)
	case *ast.BreakStatement:
		return object.Instruction(object.NativeInstructionBreak), nil
	case *ast.ContinueStatement:
		return object.Instruction(object.NativeInstructionContinue), nil
	case *ast.IdentifierExpression:
		// TODO: Guard against the possibility of getting a private value that's not in the current scope.
		value, ok := b.Get(node.Value, true)
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
		return _evalFunctionExpression(b, node)
	case *ast.CallExpression:
		return _evalCallExpression(b, node)
	default:
		return nil, runtime.NewError(
			runtime.InterpreterError,
			"I haven't been taught how to interpret this node yet!",
			runtime.ErrorValue("NodeType", fmt.Sprintf("%T", node)),
			runtime.ErrorValue("Node", node),
		)
	}
}

func _evalWhileExpression(b *binding.Binding[object.Object], node *ast.WhileExpression) (object.Object, error) {
	var result object.Object
	result = &object.Null{}

EvalLoop:
	for {
		condition, err := Evaluate(b, node.Condition)
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

		blockResult, err := _evalStatementList(b, node.Block.StatementList, true, false)
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

func _evalSubscriptNode(b *binding.Binding[object.Object], node *ast.SubscriptExpression) (object.Object, error) {
	rec, err := Evaluate(b, node.Receiver)
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

	val, err := Evaluate(b, node.Value)
	if err != nil {
		return nil, err
	}

	return sub.Subscript(val)
}

func _evalArrayNode(b *binding.Binding[object.Object], node *ast.ArrayExpression) (object.Object, error) {
	array := make([]object.Object, len(node.Elements))

	for i, n := range node.Elements {
		if v, err := Evaluate(b, n); err != nil {
			return nil, err
		} else {
			array[i] = v
		}
	}

	return object.Array(array), nil
}

func _evalNativeFunction(b *binding.Binding[object.Object], fn *object.NativeFunctionObject, node *ast.CallExpression) (object.Object, error) {
	if len(node.Arguments) != fn.ArgumentCount {
		return nil, runtime.NewError(runtime.ArgumentError,
			"Unexpected number of arguments in function call",
			runtime.ErrorValue("Expected", fn.ArgumentCount),
			runtime.ErrorValue("Received", len(node.Arguments)),
		)
	}

	arguments := make([]object.Object, len(node.Arguments))
	for i, arg := range node.Arguments {
		val, err := Evaluate(b, arg)
		if err != nil {
			return nil, err
		}
		arguments[i] = val
	}

	return fn.Body(b.NewChild(), arguments)
}

func _evalCallExpression(b *binding.Binding[object.Object], node *ast.CallExpression) (object.Object, error) {
	fn, err := Evaluate(b, node.Function)
	if err != nil {
		return nil, err
	}

	switch closure := fn.(type) {
	case *object.ClosureLiteral:
		return _evalCallClosureLiteral(b, closure, node)
	case *object.NativeFunctionObject:
		return _evalNativeFunction(b, closure, node)
	default:
		return nil, runtime.NewError(runtime.CallError,
			"Unable to call non-function type",
			runtime.ErrorValue("Function", fn),
		)
	}
}

func _evalCallClosureLiteral(b *binding.Binding[object.Object], closure *object.ClosureLiteral, node *ast.CallExpression) (object.Object, error) {
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
		val, err := Evaluate(b, arg)
		if err != nil {
			return nil, err
		}
		arguments[i] = val
	}

	child := closure.Binding.NewChild()

	for i, val := range arguments {
		name := closure.ParameterList[i]
		child.Set(name, val, binding.F_CONST|binding.F_PROTECTED)
	}

	return _evalStatementList(child, closure.Body.StatementList, true, true)
}

func _evalFunctionExpression(b *binding.Binding[object.Object], node *ast.FunctionExpression) (*object.ClosureLiteral, error) {
	parameters := collection.MapSlice(node.Parameters, func(l *ast.IdentifierExpression) string {
		return l.Value
	})
	return object.Closure(parameters, node.BlockStatement, b.NewChild()), nil
}

func _evalDeleteStatement(b *binding.Binding[object.Object], node *ast.DeleteStatement) (object.Object, error) {
	name := node.Identifier.Value

	value, currentState, _ := b.GetValueState(name, true)
	if !currentState.IsSet() {
		return nil, runtime.NewError(runtime.UnknownIdentifierError,
			"Can't delete an undefined identifier",
			runtime.ErrorValue("Location", node.SourceToken().Location),
			runtime.ErrorValue("Label", name),
		)
	}
	if !currentState.IsMutable() || currentState.IsProtected() {
		return nil, runtime.NewError(runtime.ConstantError,
			"Can't delete a constant identifier",
			runtime.ErrorValue("Location", node.SourceToken().Location),
			runtime.ErrorValue("Label", name),
		)
	}

	if !b.UnsetFirst(name, true) {
		return nil, runtime.NewError(runtime.InterpreterError,
			"Failed to delete the value for identifier",
			runtime.ErrorValue("Location", node.SourceToken().Location),
			runtime.ErrorValue("Label", name),
		)
	}

	return value, nil
}

func _evalAssignmentNode(b *binding.Binding[object.Object], node ast.AssignmentStatement, update bool) (object.Object, error) {
	name := node.Label()

	currentState := b.GetState(name, !node.CurrentFrameOnly())
	if node.RequireDefined() && !currentState.IsSet() {
		return nil, runtime.NewError(runtime.UnknownIdentifierError,
			"Unable to find identifier in scope",
			runtime.ErrorValue("Name", name),
			runtime.ErrorValue("Location", node.SourceToken().Location),
		)
	}

	if node.RequireDefined() && !currentState.IsSet() {
		return nil, runtime.NewError(runtime.UnknownIdentifierError,
			"Identifier undefined in scope",
			runtime.ErrorValue("Name", name),
			runtime.ErrorValue("Location", node.SourceToken().Location),
		)
	}

	if currentState.IsSet() && !currentState.IsMutable() {
		return nil, runtime.NewError(runtime.ConstantError,
			"Constant already defined in scope",
			runtime.ErrorValue("Name", name),
			runtime.ErrorValue("Location", node.SourceToken().Location),
		)
	}

	if currentState.IsSet() && node.RequireUndefined() {
		return nil, runtime.NewError(runtime.InterpreterError,
			"Unable to perform assignment on a defined variable",
			runtime.ErrorValue("Name", name),
			runtime.ErrorValue("Location", node.SourceToken().Location),
		)
	}

	value, err := Evaluate(b, node.AssignmentExpression())
	if err != nil {
		return nil, err
	}

	if update {
		_ = b.Update(name, value, !node.CurrentFrameOnly())
	} else {
		var f binding.Flag
		if !node.Mutable() {
			f |= binding.F_CONST
		}
		_ = b.Set(name, value, f)
	}

	return value, nil
}

func _evalReturnStatement(b *binding.Binding[object.Object], node *ast.ReturnStatement) (object.Object, error) {
	val, err := Evaluate(b, node.Expression)
	if err != nil {
		return nil, err
	}
	return object.Return(val), nil
}

func _evalIfExpression(b *binding.Binding[object.Object], node *ast.IfExpression) (object.Object, error) {
	condition, err := Evaluate(b, node.Condition)
	if err != nil {
		return nil, err
	}

	boolean := object.Bool(false)
	if err := object.CoerceToType(condition, boolean); err != nil {
		return nil, err
	}

	if boolean.Value {
		return Evaluate(b, node.TrueStatement)
	} else if node.FalseStatement != nil {
		return Evaluate(b, node.FalseStatement)
	} else {
		return &object.Null{}, nil
	}
}

func _evalInfixExpression(b *binding.Binding[object.Object], node *ast.InfixExpression) (object.Object, error) {
	lhs, err := Evaluate(b, node.Left)
	if err != nil {
		return nil, err
	}
	rhs, err := Evaluate(b, node.Right)
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

func _evalNotExpression(b *binding.Binding[object.Object], node *ast.NotExpression) (object.Object, error) {
	right, err := Evaluate(b, node.Expression)
	if err != nil {
		return nil, err
	}

	bo := new(object.Boolean)

	if err := object.CoerceToType(right, bo); err != nil {
		return nil, err
	}

	return object.Bool(!bo.Value), nil
}

func _evalNegateExpression(b *binding.Binding[object.Object], n *ast.NegateExpression) (object.Object, error) {
	obj, err := Evaluate(b, n.Expression)
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
func _evalStatementList(b *binding.Binding[object.Object], list []ast.Statement, pushStack, unwrapReturn bool) (object.Object, error) {
	if pushStack {
		b = b.NewChild()
	}
	var result object.Object
	var err error

	result = &object.Null{}

	for _, node := range list {
		result, err = Evaluate(b, node)
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
