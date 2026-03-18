package refcount

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/ebitengine/purego"
)

type State struct {
	refs atomic.Int32
}

var (
	states sync.Map
	pins   sync.Map

	addRefCallback           = purego.NewCallback(func(self unsafe.Pointer) { AddRef(self) })
	releaseCallback          = purego.NewCallback(func(self unsafe.Pointer) int32 { return Release(self) })
	hasOneRefCallback        = purego.NewCallback(func(self unsafe.Pointer) int32 { return HasOneRef(self) })
	hasAtLeastOneRefCallback = purego.NewCallback(func(self unsafe.Pointer) int32 { return HasAtLeastOneRef(self) })
)

type baseRefCounted struct {
	Size             uintptr
	AddRef           uintptr
	Release          uintptr
	HasOneRef        uintptr
	HasAtLeastOneRef uintptr
}

func Init(base unsafe.Pointer, size uintptr, owner any) *State {
	hdr := (*baseRefCounted)(base)
	hdr.Size = size
	hdr.AddRef = addRefCallback
	hdr.Release = releaseCallback
	hdr.HasOneRef = hasOneRefCallback
	hdr.HasAtLeastOneRef = hasAtLeastOneRefCallback
	state := &State{}
	state.refs.Store(1)
	key := uintptr(base)
	states.Store(key, state)
	pins.Store(key, owner)
	return state
}

func AddRef(base unsafe.Pointer) {
	if state, ok := loadState(base); ok {
		state.refs.Add(1)
	}
}

func Release(base unsafe.Pointer) int32 {
	state, ok := loadState(base)
	if !ok {
		return 1
	}
	if state.refs.Add(-1) == 0 {
		key := uintptr(base)
		states.Delete(key)
		pins.Delete(key)
		return 1
	}
	return 0
}

func HasOneRef(base unsafe.Pointer) int32 {
	if state, ok := loadState(base); ok && state.refs.Load() == 1 {
		return 1
	}
	return 0
}

func HasAtLeastOneRef(base unsafe.Pointer) int32 {
	if state, ok := loadState(base); ok && state.refs.Load() >= 1 {
		return 1
	}
	return 0
}

func Load(base unsafe.Pointer) (any, bool) {
	return pins.Load(uintptr(base))
}

func (s *State) Count() int32 { return s.refs.Load() }

func loadState(base unsafe.Pointer) (*State, bool) {
	v, ok := states.Load(uintptr(base))
	if !ok {
		return nil, false
	}
	return v.(*State), true
}
