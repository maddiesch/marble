// Package binding provides the variable binding for the evaluator
package binding

import (
	"strings"

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

// NewChild is a helper function that returns a new child of the receiving binding.
func (b *Binding[T]) NewChild() *Binding[T] {
	return New(b)
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

// UnsetFirst will delete the given key from the binding.
//
// If recursively is true, it will walk up the ancestry list until it encounters
// a binding with the given key. After it finds the key, it will stop.
//
// If recursively is false, only the receiving binding will be checked.
//
// It returns true if the value was found and deleted.
func (b *Binding[T]) UnsetFirst(key string, recursively bool) bool {
	return b.unsafeUnsetFirst(key, recursively)
}

// GetValueState will return the value and state from the binding.
//
// It's the same as calling GetState & Get
//
// The return boolean can be ignored if you're using the State IsSet function
func (b *Binding[T]) GetValueState(key string, recursively bool) (T, State, bool) {
	val, state := b.unsafeGetState(key, recursively)
	if val == nil {
		var v T
		return v, state, false
	}
	return val.Value, state, true
}

// Get the receiving binding's identifier
func (b *Binding[T]) ID() uint64 {
	return uint64(b.id)
}

// Print a debug string of the binding's current state
func (b *Binding[T]) DebugString() string {
	var builder strings.Builder

	b.unsafeDebugString(0, &builder)

	return builder.String()
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

// The value exists in the binding's ancestry
func (s State) IsSet() bool {
	return bit.Has(s, S_SET)
}

// The value exists in the current level binding, without needing to follow the
// ancestry
func (s State) IsCurrent() bool {
	return bit.Has(s, S_CURRENT_FRAME)
}

// The value is allowed to be mutated
func (s State) IsMutable() bool {
	return bit.Has(s, S_MUTABLE)
}

// The value is in a protected state
func (s State) IsProtected() bool {
	return bit.Has(s, S_PROTECTED)
}

// Returns a value's current state in the binding
func (b *Binding[T]) GetState(key string, recursively bool) State {
	_, s := b.unsafeGetState(key, recursively)
	return s
}

// Returns the binding's parent or nil if it's the root.
func (b *Binding[T]) Parent() *Binding[T] {
	return b.parent
}
