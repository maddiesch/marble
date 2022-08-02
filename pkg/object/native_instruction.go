package object

import "fmt"

const (
	NATIVE_INST ObjectType = "OBJ_NativeInstruction"
)

type NativeInstructionType uint8

const (
	_ NativeInstructionType = iota
	NativeInstructionBreak
	NativeInstructionContinue
)

func Instruction(t NativeInstructionType) *NativeInstruction {
	return &NativeInstruction{IType: t}
}

type NativeInstruction struct {
	IType NativeInstructionType
}

func (*NativeInstruction) Type() ObjectType {
	return NATIVE_INST
}

func (o *NativeInstruction) DebugString() string {
	return fmt.Sprintf("NativeInstruction(%d)", o.IType)
}

func (o *NativeInstruction) GoValue() any {
	return nil
}

func (o *NativeInstruction) CoerceTo(t ObjectType) (Object, bool) {
	return nil, false
}

var _ Object = (*NativeInstruction)(nil)
