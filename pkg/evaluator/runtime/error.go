package runtime

import (
	"fmt"
	"strings"
)

type ErrorID string

const (
	TypeError        ErrorID = "TypeError"
	InterpreterError ErrorID = "InterpreterError"
)

type Error struct {
	ID      ErrorID
	Detail  string
	Context map[string]any
}

func (r *Error) Error() string {
	var builder strings.Builder

	builder.WriteString("RuntimeError: [" + string(r.ID) + "] " + r.Detail)
	if len(r.Context) > 0 {
		builder.WriteString("\nContext:\n")
		str := make([]string, 0, len(r.Context))

		for key, value := range r.Context {
			str = append(str, fmt.Sprintf("\t%s: %s", key, value))
		}

		builder.WriteString(strings.Join(str, "\n"))
	}

	return builder.String()
}

func NewError(id ErrorID, detail string, context ...ErrorContextValue) *Error {
	ctx := make(map[string]any)

	for _, v := range context {
		ctx[v.Key()] = v.Value()
	}

	return &Error{
		ID:      id,
		Detail:  detail,
		Context: ctx,
	}
}

type ErrorContextValue interface {
	Key() string
	Value() any
}

type _errorContextValue struct {
	k string
	v any
}

func (v *_errorContextValue) Key() string {
	return v.k
}

func (v *_errorContextValue) Value() any {
	return v.v
}

func ErrorValue(k string, v any) ErrorContextValue {
	return &_errorContextValue{k, v}
}
