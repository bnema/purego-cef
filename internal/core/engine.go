// internal/core/engine.go
package core

import (
	"unsafe"

	out "github.com/bnema/purego-cef/internal/ports/out"
)

// CAPI is the outbound port interface. Core depends on this for all C
// interactions. The real implementation (internal/capi) uses purego;
// tests inject mocks.
//
// It embeds the generated AppFunctions interface for CEF free functions
// (Initialize, Shutdown, etc.) and adds infrastructure methods that are
// not in the parsed CEF headers (string ops from cef_string_types.h,
// purego callback creation, and symbol registration).
type CAPI interface {
	// Generated outbound port — CEF free functions from cef_app_capi.h
	out.AppFunctions

	// Infrastructure — not in parsed CEF headers, handwritten.
	// NewCallback wraps purego.NewCallback.
	NewCallback(fn any) uintptr
	// String operations from cef_string_types.h (not parsed by cefgen).
	StringSet(src *uint16, srcLen uintptr, output unsafe.Pointer, copy int32) int32
	StringClear(s unsafe.Pointer)
	StringUserfreeFree(s unsafe.Pointer)
	StringListSize(list uintptr) uintptr
	StringListValue(list uintptr, index uintptr, value unsafe.Pointer) int32
}

// Engine is the core domain. It holds no global state — all state
// flows through the instance. The cef/ composition root creates one
// Engine at Init() time and holds it as a package-level variable.
type Engine struct {
	capi        CAPI
	refs        *RefManager
	initialized bool
}

// New creates an Engine with the given CAPI adapter.
func New(capi CAPI) *Engine {
	e := &Engine{capi: capi}
	e.refs = NewRefManager(capi)
	return e
}

// Refs returns the engine's RefManager.
func (e *Engine) Refs() *RefManager {
	return e.refs
}
