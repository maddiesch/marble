package env

type frame struct {
	id     uint64
	lookup map[string]*Entry
}

func newFrame(id uint64) *frame {
	return &frame{
		id:     id,
		lookup: make(map[string]*Entry),
	}
}

func (f *frame) set(k string, e *Entry) {
	f.lookup[k] = e
}

func (f *frame) get(k string) (v *Entry, ok bool) {
	v, ok = f.lookup[k]
	return
}
