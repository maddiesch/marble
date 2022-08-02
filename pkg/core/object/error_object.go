package object

import (
	"fmt"
	"strconv"

	"github.com/maddiesch/marble/pkg/core/visitor"
)

const (
	ERROR ObjectType = "OBJ_Error"
)

func NewErrorf(kind, format string, args ...any) *ErrorObject {
	return NewError(kind, fmt.Sprintf(format, args...))
}

func NewError(kind, message string) *ErrorObject {
	return &ErrorObject{
		Kind:    kind,
		Message: message,
	}
}

type ErrorObject struct {
	Kind    string
	Message string
}

func (*ErrorObject) Type() ObjectType {
	return ERROR
}

func (o *ErrorObject) DebugString() string {
	return fmt.Sprintf("Error<%s>(%s)", o.Kind, strconv.Quote(o.Message))
}

func (o *ErrorObject) GoValue() any {
	return nil
}

// TODO: Delete once CastVisitor is complete
func (o *ErrorObject) CoerceTo(t ObjectType) (Object, bool) {
	switch t {
	case ERROR:
		return o, true
	default:
		return NewVoid(), false
	}
}

func (o *ErrorObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*ErrorObject)(nil)
