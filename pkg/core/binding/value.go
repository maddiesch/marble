package binding

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/debug"
)

type Option uint32

type Flag uint8

const (
	_ Flag = 1 << iota
	F_CONST
	F_PROTECTED
	F_PRIVATE
	F_NATIVE
)

type Value[T any] struct {
	Value T

	flag Flag
	id   valueID
}

func (v Value[T]) DebugString() string {
	return fmt.Sprintf("v(0x%08x, 0b%08b, %T, %s)", v.id, v.flag, v.Value, loadValueDescription(v.Value))
}

func loadValueDescription(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case debug.Description:
		return v.DebugString()
	default:
		return fmt.Sprintf("%v", v)
	}
}
