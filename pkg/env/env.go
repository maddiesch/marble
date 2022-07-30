package env

import (
	"strings"
	"sync"

	"github.com/maddiesch/marble/internal/slice"
	"github.com/maddiesch/marble/pkg/ast"
	"github.com/maddiesch/marble/pkg/native"
	"github.com/maddiesch/marble/pkg/object"
	"github.com/maddiesch/marble/pkg/stack"
	"github.com/maddiesch/marble/pkg/version"
)

type Env struct {
	mu        sync.RWMutex
	namespace *stack.Stack[string]
	lookup    []*frame
	restore   [][]*frame
	fid       uint64
	ptr       uint64
	stack     []ast.Node
}

func New() *Env {
	e := &Env{
		fid:       0,
		ptr:       1_000_000,
		namespace: stack.New[string](),
		lookup:    make([]*frame, 0, 8),
		restore:   make([][]*frame, 0, 4),
		stack:     make([]ast.Node, 0, 32),
	}

	e.PushEval(&ast.Entrypoint{})
	e.Push()

	e.set("MARBLE_VERSION", &Entry{
		Value:     object.String(version.Current),
		pointer:   1,
		protected: true,
		mutable:   false,
	})

	native.Bind(e)

	return e
}

func (e *Env) PushTo(id uint64) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	index := -1
	for i, f := range e.lookup {
		if f.id == id {
			index = i
		}
	}
	if index < 0 {
		return false
	}

	e.restore = append(e.restore, e.lookup)
	e.lookup = e.lookup[:index+1]

	return true
}

func (e *Env) Restore() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.restore) < 1 {
		panic("unable to restore without pushing back to a previous frame")
	}

	e.lookup = e.restore[len(e.restore)-1]
	e.restore = e.restore[:len(e.restore)-1]
}

func (e *Env) CurrentFrame() uint64 {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.lookup[len(e.lookup)-1].id
}

func (e *Env) PushEval(node ast.Node) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.stack = append(e.stack, node)
}

func (e *Env) PopEval() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.stack) <= 1 {
		panic("attempting to pop the final environment stack!")
	}

	e.stack = e.stack[:len(e.stack)-1]
}

func (e *Env) Push() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.fid += 1

	e.lookup = append(e.lookup, newFrame(e.fid))
}

func (e *Env) Pop() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.lookup) <= 1 {
		panic("attempting to pop the final environment lookup table!")
	}

	e.lookup = e.lookup[:len(e.lookup)-1]
}

type LabelState uint8

const (
	_ LabelState = iota
	LabelStateUnassigned
	LabelStateAssignedMutable
	LabelStateAssignedImmutable
	LabelStateAssignedProtected
)

func (s LabelState) Mutable() bool {
	return s == LabelStateUnassigned || s == LabelStateAssignedMutable
}

func (s LabelState) Defined() bool {
	return s != LabelStateUnassigned
}

func (e *Env) StateFor(key string, currentFrameOnly bool) LabelState {
	en := e.getEntry(key, !currentFrameOnly)
	if en == nil {
		return LabelStateUnassigned
	} else if en.protected {
		return LabelStateAssignedProtected
	} else if en.mutable {
		return LabelStateAssignedMutable
	} else {
		return LabelStateAssignedImmutable
	}
}

func (e *Env) SetProtected(key string, value object.Object) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.set(key, &Entry{
		Value:     value,
		mutable:   false,
		protected: false,
	})
}

func (e *Env) Set(key string, value object.Object, mutable bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.set(key, &Entry{
		Value:   value,
		mutable: mutable,
	})
}

func (e *Env) set(key string, entry *Entry) {
	e.ptr += 1

	frame := e.lookup[len(e.lookup)-1]
	entry.frame = frame.id
	entry.pointer = e.ptr
	frame.set(key, entry)
}

func (e *Env) GetEntry(key string) *Entry {
	return e.getEntry(key, true)
}

func (e *Env) getEntry(key string, recursively bool) *Entry {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for i := len(e.lookup) - 1; i >= 0; i-- {
		frame := e.lookup[i]

		if e, ok := frame.get(key); ok {
			return e
		}

		if !recursively {
			return nil
		}
	}

	return nil
}

func (e *Env) Delete(key string, currentFrameOnly bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i := len(e.lookup) - 1; i >= 0; i-- {
		if e.lookup[i].delete(key) {
			return
		}

		if currentFrameOnly {
			return
		}
	}
}

func (e *Env) Get(key string) (object.Object, bool) {
	entry := e.GetEntry(key)
	if entry == nil {
		return &object.Void{}, false
	}
	return entry.Value, true
}

func (e *Env) DebugString() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var builder strings.Builder

	builder.WriteString("Environment\n")
	frames := slice.Map(e.lookup, func(f *frame) string {
		return f.debugString(1)
	})

	builder.WriteString(strings.Join(frames, "\n"))

	return builder.String()
}

func (e *Env) PushNS(string) {
}

func (e *Env) PopNS() {
}

var _ object.Binding = (*Env)(nil)
