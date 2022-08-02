package visitor

type Visitor[T any] interface {
	Visit(T)
}
