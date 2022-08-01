package binding

import (
	"sync/atomic"

	"github.com/maddiesch/marble/internal/bit"
)

type valueID uint64
type frameID uint64

var (
	_valueID uint64
	_frameID uint64
)

func getNextValID() valueID {
	return valueID(atomic.AddUint64(&_valueID, 1))
}

func getNextFrameID() frameID {
	return frameID(atomic.AddUint64(&_frameID, 1))
}

func (b *Binding[T]) unsafeSet(k string, v T, f Flag) bool {
	if _, ok, _ := b.unsafeGet(k, false, 0); !ok {
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
		return nil, false, -1
	}
}

func (b *Binding[T]) unsafeGetState(k string, r bool) State {
	val, ok, d := b.unsafeGet(k, r, 0)
	if !ok {
		return 0
	}
	s := S_SET

	if d == 0 {
		s |= S_CURRENT_FRAME
	}
	if !bit.Has(val.flag, F_CONST) {
		s |= S_MUTABLE
	}
	if !bit.Has(val.flag, F_PRIVATE) {
		s |= S_PRIVATE
	}
	if !bit.Has(val.flag, F_PROTECTED) {
		s |= S_PROTECTED
	}

	return s
}
