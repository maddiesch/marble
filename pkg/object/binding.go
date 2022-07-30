package object

import "io"

type Binding interface {
	SetProtected(string, Object)
	Get(string) (Object, bool)
	Set(string, Object, bool)
	Update(string, Object, bool) bool
	PushNS(string)
	PopNS()
	Stdout() io.Writer
	Stderr() io.Writer
}
