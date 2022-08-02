package object

const (
	NULL ObjectType = "OBJ_Null"
)

type Null struct {
	Value bool
}

func (o *Null) Type() ObjectType {
	return NULL
}

func (o *Null) DebugString() string {
	return "Null()"
}

func (o *Null) GoValue() any {
	return nil
}

func (o *Null) CoerceTo(t ObjectType) (Object, bool) {
	switch t {
	case NULL:
		return o, true
	case BOOLEAN:
		return Bool(false), true
	default:
		return &Void{}, false
	}
}

var _ Object = (*Null)(nil)
