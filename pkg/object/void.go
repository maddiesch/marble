package object

const (
	VOID ObjectType = "OBJ_Void"
)

type Void struct {
	Value bool
}

func (*Void) Type() ObjectType {
	return VOID
}

func (*Void) Description() string {
	return "VOID"
}

func (*Void) GoValue() any {
	return nil
}

func (o *Void) CoerceTo(t ObjectType) (Object, bool) {
	if t == o.Type() {
		return o, true
	}
	return &Void{}, false
}

var _ Object = (*Void)(nil)
