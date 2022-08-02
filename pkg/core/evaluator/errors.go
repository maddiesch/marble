package evaluator

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/ast"
	"github.com/maddiesch/marble/pkg/core/object"
)

type IllegalExpression struct {
	Node    ast.Node
	Message string
}

func (e IllegalExpression) Error() string {
	return fmt.Sprintf("IllegalExpression: %s\n\tSource: %s", e.Message, e.Node.SourceToken().Location)
}

type ObjectCastError struct {
	Node   ast.Node
	Object object.Object
	Target object.ObjectType
}

func (e ObjectCastError) Error() string {
	return fmt.Sprintf("ObjectCastError: Unable to cast from %s to %s\n\tSource: %s", e.Object.Type(), e.Target, e.Node.SourceToken().Location)
}

type TypeError struct {
}
