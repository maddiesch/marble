package binding

import (
	"fmt"
	"io"
	"strings"
	"sync/atomic"

	"github.com/maddiesch/marble/internal/bit"
	"github.com/maddiesch/marble/internal/collection"
)

type valueID uint64
type frameID uint64

var (
	_valueID uint64 = 0x0000_1000
	_frameID uint64 = 0x0001_0000
)

func getNextValID() valueID {
	return valueID(atomic.AddUint64(&_valueID, 1))
}

func getNextFrameID() frameID {
	return frameID(atomic.AddUint64(&_frameID, 1))
}

func (b *Binding[T]) unsafeSet(k string, v T, f Flag) bool {
	if _, ok, _ := b.unsafeGet(k, false, 0); ok {
		return false
	}
	b.table[k] = &Value[T]{
		id:    getNextValID(),
		flag:  f,
		Value: v,
	}
	return true
}

func (b *Binding[T]) unsafeUpdate(k string, v T, r bool) bool {
	val, ok, _ := b.unsafeGet(k, r, 0)
	if !ok {
		return false
	}
	val.Value = v

	return true
}

func (b *Binding[T]) unsafeGet(k string, r bool, d int) (*Value[T], bool, int) {
	if v, ok := b.table[k]; ok {
		return v, true, d
	} else if r && b.parent != nil {
		return b.parent.unsafeGet(k, r, d+1)
	} else {
		return nil, false, -(d + 1)
	}
}

func (b *Binding[T]) unsafeGetState(k string, r bool) (*Value[T], State) {
	val, ok, d := b.unsafeGet(k, r, 0)
	if !ok {
		return nil, 0
	}
	s := S_SET

	if d == 0 {
		s |= S_CURRENT_FRAME
	}
	if !bit.Has(val.flag, F_CONST) {
		s |= S_MUTABLE
	}
	if bit.Has(val.flag, F_PRIVATE) {
		s |= S_PRIVATE
	}
	if bit.Has(val.flag, F_PROTECTED) {
		s |= S_PROTECTED
	}

	return val, s
}

func (b *Binding[T]) unsafeUnsetFirst(k string, r bool) bool {
	if _, ok := b.table[k]; ok {
		delete(b.table, k)
		return true
	} else if r && b.parent != nil {
		return b.parent.unsafeUnsetFirst(k, r)
	} else {
		return false
	}
}

func (b *Binding[T]) unsafeDebugString(d int, w io.Writer) {
	prefix := strings.Repeat("\t", d)

	fmt.Fprintf(w, "%sBinding (0x%08x)\n", prefix, b.id)

	for _, key := range collection.SortedKeys(b.table) {
		fmt.Fprintf(w, "%s\t%s = %s\n", prefix, key, b.table[key].DebugString())
	}

	if b.parent != nil {
		b.parent.unsafeDebugString(d+2, w)
	}
}
