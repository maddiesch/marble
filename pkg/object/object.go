package object

type Object interface {
	Type() ObjectType

	GoValue() any

	Description() string

	Cast(ObjectType) (Object, bool)
}

func CastChain(o Object, t ...ObjectType) (Object, bool) {
	r := o
	var ok bool
	for _, t := range t {
		r, ok = r.Cast(t)
		if !ok {
			return nil, false
		}
	}
	return r, true
}
