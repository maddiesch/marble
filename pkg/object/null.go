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

func (o *Null) Description() string {
	return "NULL"
}

func (o *Null) GoValue() any {
	return nil
}

func (o *Null) Cast(t ObjectType) (Object, bool) {
	if t == o.Type() {
		return o, true
	}
	return &Void{}, false
}

var _ Object = (*Null)(nil)
