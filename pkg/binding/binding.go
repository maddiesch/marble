// Package binding provides the variable binding for the evaluator
package binding

import (
	"github.com/maddiesch/marble/internal/bit"
	"github.com/maddiesch/marble/pkg/debug"
)

// Binding contains a table to map keys to values as well as a hierarchy of
// parent bindings
type Binding[T any] struct {
	debug.Description

	id     frameID
	parent *Binding[T]

	table map[string]*Value[T]
}

// New returns a new Binding as a child of the parent binding.
//
// Note: nil is a valid parent binding
func New[T any](parent *Binding[T]) *Binding[T] {
	return &Binding[T]{
		id:     getNextFrameID(),
		parent: parent,
		table:  make(map[string]*Value[T]),
	}
}

// Set will update a value in the receiving context
func (b *Binding[T]) Set(key string, value T, flag Flag) bool {
	return b.unsafeSet(key, value, flag)
}

// Update will update an existing value with the option to recursively walk the
// ancestry to find the value to update.
//
// If no value is found to be updated, the return value will be false
func (b *Binding[T]) Update(key string, value T, recursively bool) bool {
	return b.unsafeUpdate(key, value, recursively)
}

// Get will search the mapping table for the given key, and optionally walk the
// ancestry until the value is found. If no value is found, the return value
// will be the zero value for the binding's containing type and the boolean will
// be false
func (b *Binding[T]) Get(key string, recursively bool) (T, bool) {
	if val, ok, _ := b.unsafeGet(key, recursively, 0); ok {
		return val.Value, true
	} else {
		var v T
		return v, false
	}
}

// Print a debug string of the binding's current state
func (b *Binding[T]) DebugString() string {
	panic("Binding.DebugString")
}

// State represents the current state of a value in the binding
type State uint8

const (
	_               State = 1 << iota
	S_SET                 // The value has been set
	S_MUTABLE             // The value is mutable
	S_PROTECTED           // The value is protected
	S_PRIVATE             // The value is private
	S_CURRENT_FRAME       // The value was found in the first level binding
)

// The value existing the binding's ancestry
func (s State) IsSet() bool {
	return bit.Has(s, S_SET)
}

// The value exists in the current level binding, without needing to follow the ancestry
func (s State) IsCurrent() bool {
	return bit.Has(s, S_CURRENT_FRAME)
}

// IsMutable the value is allowed to be mutated
func (s State) IsMutable() bool {
	return bit.Has(s, S_MUTABLE)
}

// Returns a value's current state in the binding
func (b *Binding[T]) GetState(key string, recursively bool) State {
	return b.unsafeGetState(key, recursively)
}
