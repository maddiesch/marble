package env

import (
	"fmt"
	"strings"

	"github.com/maddiesch/marble/internal/collection"
)

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

func (f *frame) delete(k string) bool {
	if _, ok := f.lookup[k]; ok {
		delete(f.lookup, k)
		return true
	}
	return false
}

func (f *frame) debugString(ident int) string {
	var builder strings.Builder
	builder.WriteString(strings.Repeat("\t", ident))
	builder.WriteString(fmt.Sprintf("Frame (%d)", f.id))

	entries := collection.MapMap(f.lookup, func(k string, v *Entry) string {
		return fmt.Sprintf("%s%s = %s", strings.Repeat("\t", ident+1), k, v.Value.Description())
	})

	if len(entries) > 0 {
		builder.WriteString("\n")
	}

	builder.WriteString(strings.Join(entries, "\n"))

	return builder.String()
}
