package env

import (
	"sync"

	"github.com/maddiesch/marble/pkg/ast"
	"github.com/maddiesch/marble/pkg/object"
	"github.com/maddiesch/marble/pkg/version"
)

type Env struct {
	mu     sync.RWMutex
	lookup []*frame
	fid    uint64
	ptr    uint64
	stack  []ast.Node
}

func New() *Env {
	e := &Env{
		fid:    0,
		ptr:    1_000_000,
		lookup: []*frame{},
		stack:  make([]ast.Node, 0, 32),
	}

	e.PushEval(&ast.Entrypoint{})
	e.Push()

	e.set("MARBLE_VERSION", &Entry{
		Value:     object.String(version.Current),
		pointer:   1,
		protected: true,
		mutable:   false,
	})

	return e
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

func (e *Env) Set(key string, value object.Object, mutable bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.ptr += 1

	e.set(key, &Entry{
		Value:   value,
		pointer: e.ptr,
		mutable: mutable,
	})
}

func (e *Env) set(key string, entry *Entry) {
	frame := e.lookup[len(e.lookup)-1]
	entry.frame = frame.id
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

func (e *Env) Get(key string) (object.Object, bool) {
	entry := e.GetEntry(key)
	if entry == nil {
		return &object.Void{}, false
	}
	return entry.Value, true
}
