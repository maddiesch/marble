package binding

type Option uint32

type Flag uint8

const (
	_ Flag = 1 << iota
	F_CONST
	F_PROTECTED
	F_PRIVATE
)

type Value[T any] struct {
	Value T

	flag Flag
	id   valueID
}
