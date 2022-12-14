package object

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/core/visitor"
)

const (
	NATIVE_INST ObjectType = "OBJ_NativeInstruction"
)

type NativeInstructionType uint8

const (
	_ NativeInstructionType = iota
	NativeInstructionBreak
	NativeInstructionContinue
)

func NewNativeInstruction(t NativeInstructionType) *NativeInstructionObject {
	return &NativeInstructionObject{IType: t}
}

type NativeInstructionObject struct {
	IType NativeInstructionType
}

func (*NativeInstructionObject) Type() ObjectType {
	return NATIVE_INST
}

func (o *NativeInstructionObject) DebugString() string {
	return fmt.Sprintf("NativeInstruction(%d)", o.IType)
}

func (o *NativeInstructionObject) GoValue() any {
	return nil
}

func (o *NativeInstructionObject) Accept(v visitor.Visitor[Object]) {
	v.Visit(o)
}

var _ Object = (*NativeInstructionObject)(nil)
