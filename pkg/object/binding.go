package object

type Binding interface {
	SetProtected(string, Object)
	Get(string) (Object, bool)
	Set(string, Object, bool)
	Update(string, Object, bool) bool
	PushNS(string)
	PopNS()
}
