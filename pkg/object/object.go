package object

type Object interface {
	Type() ObjectType

	GoValue() any

	Description() string

	Cast(ObjectType) (Object, bool)
}
