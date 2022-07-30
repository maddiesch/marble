package object

type Binding interface {
	SetProtected(string, Object)
	Get(string) (Object, bool)
	Set(string, Object, bool)
	PushNS(string)
	PopNS()
}
