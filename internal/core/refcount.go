// internal/core/refcount.go
package core

import (
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

// refState tracks the atomic reference count for a single Go-owned CEF object.
type refState struct {
	refs   atomic.Int32
	pinner runtime.Pinner
}

// baseRefCounted mirrors the cef_base_ref_counted_t layout.
type baseRefCounted struct {
	Size             uintptr
	AddRef           uintptr
	Release          uintptr
	HasOneRef        uintptr
	HasAtLeastOneRef uintptr
}

// RefManager handles CEF reference counting for Go-owned objects.
type RefManager struct {
	states sync.Map // uintptr -> *refState
	pins   sync.Map // uintptr -> any (the owner, kept alive)

	// Pre-created callbacks (created once via CAPI.NewCallback)
	addRefCb           uintptr
	releaseCb          uintptr
	hasOneRefCb        uintptr
	hasAtLeastOneRefCb uintptr
}

// NewRefManager creates a RefManager, registering callbacks via the CAPI adapter.
func NewRefManager(capi CAPI) *RefManager {
	rm := &RefManager{}
	rm.addRefCb = capi.NewCallback(func(self unsafe.Pointer) { rm.addRef(self) })
	rm.releaseCb = capi.NewCallback(func(self unsafe.Pointer) int32 { return rm.release(self) })
	rm.hasOneRefCb = capi.NewCallback(func(self unsafe.Pointer) int32 { return rm.hasOneRef(self) })
	rm.hasAtLeastOneRefCb = capi.NewCallback(func(self unsafe.Pointer) int32 { return rm.hasAtLeastOneRef(self) })
	return rm
}

// InitRefCount wires refcount callbacks into the cef_base_ref_counted_t header.
// size is unsafe.Sizeof of the full C struct. owner is pinned to prevent GC relocation.
func (rm *RefManager) InitRefCount(base unsafe.Pointer, size uintptr, owner any) {
	if base == nil {
		panic("InitRefCount: nil base pointer")
	}
	hdr := (*baseRefCounted)(base)
	hdr.Size = size
	hdr.AddRef = rm.addRefCb
	hdr.Release = rm.releaseCb
	hdr.HasOneRef = rm.hasOneRefCb
	hdr.HasAtLeastOneRef = rm.hasAtLeastOneRefCb

	state := &refState{}
	state.refs.Store(1)
	state.pinner.Pin(owner)

	key := uintptr(base)
	rm.states.Store(key, state)
	rm.pins.Store(key, owner)
}

func (rm *RefManager) addRef(base unsafe.Pointer) {
	if st, ok := rm.loadState(base); ok {
		st.refs.Add(1)
	}
}

func (rm *RefManager) release(base unsafe.Pointer) int32 {
	st, ok := rm.loadState(base)
	if !ok {
		return 1
	}
	refs := st.refs.Add(-1)
	if refs == 0 {
		key := uintptr(base)
		rm.states.Delete(key)
		rm.pins.Delete(key)
		st.pinner.Unpin()
		return 1
	}
	return 0
}

func (rm *RefManager) hasOneRef(base unsafe.Pointer) int32 {
	if st, ok := rm.loadState(base); ok && st.refs.Load() == 1 {
		return 1
	}
	return 0
}

func (rm *RefManager) hasAtLeastOneRef(base unsafe.Pointer) int32 {
	if st, ok := rm.loadState(base); ok && st.refs.Load() >= 1 {
		return 1
	}
	return 0
}

// AddRef increments the refcount for the object at base.
func (rm *RefManager) AddRef(base unsafe.Pointer) {
	rm.addRef(base)
}

// LoadOwner retrieves the Go owner struct stored for the given base pointer.
func (rm *RefManager) LoadOwner(base unsafe.Pointer) (any, bool) {
	return rm.pins.Load(uintptr(base))
}

func (rm *RefManager) loadState(base unsafe.Pointer) (*refState, bool) {
	v, ok := rm.states.Load(uintptr(base))
	if !ok {
		return nil, false
	}
	return v.(*refState), true
}
