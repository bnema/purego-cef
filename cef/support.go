package cef

import (
	"runtime"
	"sync"
	"sync/atomic"
	"unicode/utf16"
	"unsafe"

	"github.com/bnema/purego-cef/cef/internal/raw"
	"github.com/ebitengine/purego"
)

// ---------------------------------------------------------------------------
// String conversion (ported from internal/cefstr)
// ---------------------------------------------------------------------------

var (
	stringSet    func(*uint16, uintptr, *raw.CEFStringT, int32) int32
	stringClear  func(*raw.CEFStringT)
	stringFreeFn func(*raw.CEFStringT)
	stringsBound atomic.Bool
)

// bindStringFuncs registers the CEF UTF-16 string functions from the loaded
// library handle.
func bindStringFuncs(handle uintptr) {
	purego.RegisterLibFunc(&stringSet, handle, "cef_string_utf16_set")
	purego.RegisterLibFunc(&stringClear, handle, "cef_string_utf16_clear")
	purego.RegisterLibFunc(&stringFreeFn, handle, "cef_string_userfree_utf16_free")
	stringsBound.Store(true)
}

// cefString converts a Go string to a CEF UTF-16 string. The caller must
// call freeCefString on the result when done.
func cefString(s string) raw.CEFStringT {
	encoded := utf16.Encode([]rune(s))
	var src *uint16
	if len(encoded) > 0 {
		src = &encoded[0]
	}
	var out raw.CEFStringT
	stringSet(src, uintptr(len(encoded)), &out, 1)
	runtime.KeepAlive(encoded)
	return out
}

// freeCefString releases a CEF string's backing memory.
func freeCefString(cs *raw.CEFStringT) {
	if cs != nil {
		stringClear(cs)
	}
}

// goString converts a pointer to a CEF string (cef_string_t*) to a Go string.
func goString(cs unsafe.Pointer) string {
	if cs == nil {
		return ""
	}
	v := (*raw.CEFStringT)(cs)
	if v.Str == nil || v.Length == 0 {
		return ""
	}
	slice := unsafe.Slice(v.Str, v.Length)
	return string(utf16.Decode(slice))
}

// goStringUserfree converts a cef_string_userfree_t to a Go string and frees
// the userfree string. The pointer parameter is unsafe.Pointer because CEF
// returns userfree strings as opaque handles.
func goStringUserfree(ptr unsafe.Pointer) string {
	if ptr == nil {
		return ""
	}
	cs := (*raw.CEFStringT)(ptr)
	s := goString(ptr)
	stringFreeFn(cs)
	return s
}

// ---------------------------------------------------------------------------
// Refcount management (ported from internal/refcount)
// ---------------------------------------------------------------------------

// refState tracks the atomic reference count for a single Go-owned CEF object.
type refState struct {
	refs   atomic.Int32
	pinner runtime.Pinner
}

var (
	refStates sync.Map // uintptr -> *refState
	refPins   sync.Map // uintptr -> any (the owner, kept alive)

	// Pre-created callbacks so purego.NewCallback is called only once.
	addRefCb           = purego.NewCallback(func(self unsafe.Pointer) { addRef(self) })
	releaseCb          = purego.NewCallback(func(self unsafe.Pointer) int32 { return release(self) })
	hasOneRefCb        = purego.NewCallback(func(self unsafe.Pointer) int32 { return hasOneRef(self) })
	hasAtLeastOneRefCb = purego.NewCallback(func(self unsafe.Pointer) int32 { return hasAtLeastOneRef(self) })
)

// baseRefCounted mirrors the cef_base_ref_counted_t layout.
type baseRefCounted struct {
	Size             uintptr
	AddRef           uintptr
	Release          uintptr
	HasOneRef        uintptr
	HasAtLeastOneRef uintptr
}

// initRefCount wires refcount callbacks into the cef_base_ref_counted_t header
// at base. size is unsafe.Sizeof of the full C struct. owner is the Go struct
// that backs the C memory -- it is pinned so the GC cannot move it.
func initRefCount(base unsafe.Pointer, size uintptr, owner any) {
	if base == nil {
		panic("initRefCount: nil base pointer")
	}
	hdr := (*baseRefCounted)(base)
	hdr.Size = size
	hdr.AddRef = addRefCb
	hdr.Release = releaseCb
	hdr.HasOneRef = hasOneRefCb
	hdr.HasAtLeastOneRef = hasAtLeastOneRefCb

	state := &refState{}
	state.refs.Store(1)
	state.pinner.Pin(owner)

	key := uintptr(base)
	refStates.Store(key, state)
	refPins.Store(key, owner)
}

func addRef(base unsafe.Pointer) {
	if st, ok := loadRefState(base); ok {
		st.refs.Add(1)
	}
}

func release(base unsafe.Pointer) int32 {
	st, ok := loadRefState(base)
	if !ok {
		return 1
	}
	if st.refs.Add(-1) == 0 {
		key := uintptr(base)
		refStates.Delete(key)
		refPins.Delete(key)
		st.pinner.Unpin()
		return 1
	}
	return 0
}

func hasOneRef(base unsafe.Pointer) int32 {
	if st, ok := loadRefState(base); ok && st.refs.Load() == 1 {
		return 1
	}
	return 0
}

func hasAtLeastOneRef(base unsafe.Pointer) int32 {
	if st, ok := loadRefState(base); ok && st.refs.Load() >= 1 {
		return 1
	}
	return 0
}

// loadRefOwner retrieves the Go owner struct stored for the given base pointer.
func loadRefOwner(base unsafe.Pointer) (any, bool) {
	return refPins.Load(uintptr(base))
}

// rawPointerHolder is implemented by all generated xxxImpl types that wrap
// a raw CEF struct pointer.  It allows handler callbacks to extract the
// underlying pointer when the user returns an interface value.
type rawPointerHolder interface {
	rawPointer() unsafe.Pointer
}

// extractRawPointer returns the underlying raw CEF pointer from an interface
// value, or nil if the value is nil or does not implement rawPointerHolder.
func extractRawPointer(v any) unsafe.Pointer {
	if v == nil {
		return nil
	}
	if h, ok := v.(rawPointerHolder); ok {
		return h.rawPointer()
	}
	return nil
}

func loadRefState(base unsafe.Pointer) (*refState, bool) {
	v, ok := refStates.Load(uintptr(base))
	if !ok {
		return nil, false
	}
	return v.(*refState), true
}
